package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/cribb/gator/internal/config"
	"github.com/cribb/gator/internal/database"
	"github.com/google/uuid"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	// name        string
	// description string
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {

	f, ok := c.handlers[cmd.name]
	if !ok {
		return fmt.Errorf("Command %s not found.", cmd.name)
	}

	err := f(s, cmd)
	if err != nil {
		return err
	}
	return nil
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) == 0 {
		return fmt.Errorf("Error: no username provided.")
	}

	providedUserName := cmd.args[0]

	user, err := s.db.GetUser(context.Background(), providedUserName)
	if err != nil {
		return fmt.Errorf("Error when logging in user '%v': (%w)", providedUserName, err)
	}

	// fmt.Printf("cmd.args: %v\n", cmd.args)
	// user := cmd.args[0]
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		// fmt.Println("failed to set user.")
		return err
	}
	fmt.Printf("User set to: %v\n", user.Name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("Error: no username provided.")
	} else if len(cmd.args) > 1 {
		return fmt.Errorf("Error: too many parameters.")
	}

	newUserName := cmd.args[0]
	now := time.Now().UTC()

	params := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      newUserName,
	}
	// fmt.Printf("params: %v\n", params)

	user, err := s.db.CreateUser(context.Background(), params)
	if err != nil {
		return fmt.Errorf("Error: failed when creating database user '%v'.\n  (%w)", newUserName, err)
	}
	fmt.Printf("created user '%v'\n", user.Name)
	// fmt.Printf("userdata: %v\n", user)

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		// fmt.Println("failed to set user.")
		return err
	}
	fmt.Printf("User set to: %v\n", user.Name)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("Error generating user list.\n  (%w)\n", err)
	}

	currentUser := s.cfg.CurrentUsername

	fmt.Println("List of Users:")
	for _, user := range users {
		fmt.Printf(" * %s", user)
		if user == currentUser {
			fmt.Printf(" (current)")
		}
		fmt.Println()
	}
	fmt.Println()

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("Error generating feed list.\n  (%w)\n", err)
	}

	fmt.Println("List of Feeds:")
	for _, feed := range feeds {
		fmt.Printf(" * %s\n", feed.FeedName)
	}
	fmt.Println()

	// fmt.Printf("feeds: \n %+v\n", feeds)
	return nil

}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetDatabase(context.Background())
	if err != nil {
		return fmt.Errorf("Error reseting feed_follows table.\n  (%w)\n", err)
	}
	fmt.Println("Resetting database: successful")

	return nil
}

func handlerAgg(s *state, cmd command) error {
	N := len(cmd.args)
	if N != 1 {
		return fmt.Errorf("Error: incorrect numer of arguments.\n Usage: addfeed <interval>")
	}

	intervalString := cmd.args[0]

	interval, err := time.ParseDuration(intervalString)
	if err != nil {
		return err
	}

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()

	fmt.Printf("Collecting feeds every %v\n", interval)
	ticker := time.NewTicker(interval)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}

	// feed, err := fetchFeed(ctx, feedURL)
	// if err != nil {
	// 	return fmt.Errorf("Feed not fetched. Well's run dry.")
	// }

	// fmt.Printf("feed: %v \n", feed)
	// return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	N := len(cmd.args)
	if N != 2 {
		return fmt.Errorf("Error: incorrect numer of arguments.\n Usage: addfeed <name> <url>")
	}

	feedName := cmd.args[0]
	feedUrl := cmd.args[1]

	// user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)
	// if err != nil {
	// 	return err
	// }

	now := time.Now().UTC()

	params := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	}

	feed, err := s.db.CreateFeed(context.Background(), params)
	if err != nil {
		return err
	}
	fmt.Printf("Added feed...\n  * %+v\t(%+v)\n", feed.Name, feed.Url)

	followParams := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedRow, err := s.db.CreateFeedFollow(context.Background(), followParams)
	if err != nil {
		return err
	}

	fmt.Printf("User: %v\n", feedRow.UserName)
	fmt.Printf("Followed feed: %v\n", feedRow.FeedName)

	return nil
}

func handlerFollow(s *state, cmd command, user database.User) error {
	N := len(cmd.args)
	if N != 1 {
		return fmt.Errorf("Error: incorrect numer of arguments.\n Usage: follow <url>")
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	// currentUser := s.cfg.CurrentUsername
	// user, err := s.db.GetUser(context.Background(), currentUser)
	// if err != nil {
	// 	return err
	// }

	// fmt.Printf("feed ID: %v\n", feed.ID)
	// fmt.Printf("User ID: %v\n", user.ID)

	now := time.Now().UTC()

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	}

	feedRow, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	// fmt.Printf("feedRow: %+v\n", feedRow)
	fmt.Printf("User: %v\n", feedRow.UserName)
	fmt.Printf("Followed feed: %v\n", feedRow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {

	// currentUser := s.cfg.CurrentUsername
	// user, err := s.db.GetUser(context.Background(), currentUser)
	// if err != nil {
	// 	return err
	// }

	feeds, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return err
	}

	fmt.Printf("User: %v\n", user.Name)
	// fmt.Printf("Feeds: %v\n", feeds)

	fmt.Printf("Following: \n")
	for _, feed := range feeds {
		fmt.Printf(" * %s\n", feed.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	// find userid from user name -- actually, no, because the user id is passed in as cmd.argument

	N := len(cmd.args)
	if N != 1 {
		return fmt.Errorf("Error: incorrect numer of arguments.\n Usage: unfollow <url>")
	}

	url := cmd.args[0]

	feed, err := s.db.GetFeedByUrl(context.Background(), url)
	if err != nil {
		return err
	}

	params := database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	}

	err = s.db.DeleteFeedFollow(context.Background(), params)
	if err != nil {
		return err
	}

	fmt.Printf("User %v no longer follows \"%v\" (%v).\n", user.Name, feed.Name, feed.Url)
	return nil
}

func handlerBrowse(s *state, cmd command, user database.User) error {
	N := len(cmd.args)

	var limit int32

	if N > 1 {
		return fmt.Errorf("Error: incorrect numer of arguments.\n Usage: browse <limit>")
	}

	if N == 1 {
		limit64, err := strconv.ParseInt(cmd.args[0], 10, 32)
		if err != nil {
			return err
		}
		limit = int32(limit64)
	}

	if N < 1 {
		fmt.Println("Setting limit to default (2).")
		limit = 2
	}

	fmt.Printf("Printing %v posts to the console...\n", limit)

	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	}

	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		return err
	}

	for i, post := range posts {
		fmt.Printf("  %v. %v\n", i+1, post.Title)
	}

	return nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {

	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUsername)

		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}

}

// func (c *commands) handlerHelp(s *state, cmd command) error {
// 	fmt.Printf("\nWelcome to Gator, the Aggregator!\nUsage:\n\n")
// 	commands := getCommands()

// 	for _, command := range commands {
// 		fmt.Printf(" %v\t%v\n", command.name, command.description)
// 	}
// 	fmt.Println()
// 	return nil
// }
