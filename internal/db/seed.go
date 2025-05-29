package db

import (
	"context"

	"database/sql"
	"fmt"
	"math/rand"

	"github.com/LikhithMar14/gopher-chat/internal/models"
	"github.com/LikhithMar14/gopher-chat/internal/service"
	"github.com/LikhithMar14/gopher-chat/internal/utils"
	"go.uber.org/zap"
)

var Usernames = []string{
	"raj_kumar", "priya_sharma", "amit_patel", "sunita_verma", "rahul_mehta",
	"ananya_singh", "vivek_yadav", "neha_rani", "arjun_das", "isha_gupta",
	"manish_chopra", "rekha_nair", "suresh_babu", "divya_jain", "kiran_thakur",
	"nilesh_mishra", "sonali_giri", "ravi_saxena", "tanya_pandey", "yash_dhawan",
	"harshita_iyer", "vinit_joseph", "megha_reddy", "deepak_nambiar", "swati_pillai",
	"rohit_sen", "alka_kapoor", "sachin_khatri", "pooja_goswami", "kunal_rao",
	"radhika_patel", "naveen_chaturvedi", "tanvi_bisht", "jay_mukherjee", "anita_kumar",
	"vishal_tiwari", "lavanya_naidu", "omkar_singhania", "shruti_menon", "abhishek_pillai",
	"meera_saxena", "sanjay_kashyap", "aarti_dubey", "dev_singh", "bhavna_shetty",
	"tejas_bansal", "niharika_rai", "gopal_malhotra", "trisha_jose", "karan_nayak",
}

var Titles = []string{
	"Exploring Indian Street Food", "Monsoon Memories", "My Solo Trip to Manali",
	"How I Cleared GATE", "Delhi Metro Diaries", "Bollywood vs South Movies",
	"Startup Culture in Bangalore", "Best Trekking Spots in India", "Tea vs Coffee in India",
	"How to Crack UPSC", "Traditional South Indian Breakfasts", "Top Engineering Colleges in India",
	"Is Coding Worth It?", "Kerala Backwaters Review", "Lassi vs Buttermilk",
	"Favourite IPL Moments", "Dussehra in Mysore", "Railway Journeys in India",
	"College Life in Tier-3 Colleges", "How I Got My First Job", "Best Indian Web Series",
	"Holi Celebration in Varanasi", "My Journey Learning Go", "How to Survive in Mumbai",
	"Best Budget Phones in India", "Train Travel Tips", "Delhi Street Shopping",
	"Life of a CS Student", "Top 10 Indian Dishes", "Indian Wedding Functions",
	"Best Time to Visit Goa", "Biryani from Hyderabad", "Indian Parenting Styles",
	"Diwali Preparation Tips", "How to Start Freelancing", "Are Indian Universities Outdated?",
	"Temple Architecture in South India", "Learning from Indian Mythology",
	"Exam Stress Relievers", "Why Indian Railways is Fascinating", "Cycling in Indian Cities",
	"Pollution Problem in Delhi", "Top Indian Instagram Influencers", "Trekking in Himachal",
	"Online Dating in India", "Food Challenges with Friends", "Study Abroad from India",
	"Growing as a Developer in India", "Yoga and Meditation Benefits", "College Fest Memories",
}

var Contents = []string{
	"This post is about my recent visit to Chandni Chowk where I explored chaat, jalebi, and more.",
	"Reminiscing the monsoon rains in Mumbai and the joy of chai-pakoda evenings.",
	"My 7-day solo adventure to Manali – food, views, and solo vibes!",
	"Sharing tips and resources I used to crack GATE in my 3rd attempt.",
	"A funny experience I had while travelling on Delhi metro during peak hours.",
	"My comparison of Bollywood and South Indian movies – both unique in their own way.",
	"What it's really like to work at a startup in Bangalore.",
	"A guide to best trekking destinations in India for adventure lovers.",
	"Discussing the eternal debate – tea vs coffee among Indians.",
	"A breakdown of UPSC prep strategy that worked for me.",
	"An overview of delicious breakfast items from Tamil Nadu and Kerala.",
	"My experience visiting various engineering colleges before admission.",
	"Discussing whether coding is a good career in India.",
	"Review of my trip to Kerala and the houseboat stays.",
	"Comparing lassi and buttermilk – what's your favorite?",
	"Sharing some iconic IPL moments that I'll never forget.",
	"My trip to Mysore during Dussehra – cultural festivity at its best.",
	"Travelling 32 hours in Indian trains – here's what I learned.",
	"How it feels to be a student in a Tier-3 Indian college.",
	"My job hunt story and how I got placed in an MNC.",
	"A review of Indian web series that kept me binge-watching.",
	"Experience of celebrating Holi in Varanasi – colorful chaos!",
	"Sharing my experience of learning Go programming from scratch.",
	"Tips on how to survive in Mumbai as a newcomer.",
	"A comparison of budget phones under 20k in India.",
	"Some useful tips to travel long distances by train in India.",
	"Tips for budget street shopping in Delhi markets.",
	"My daily routine as a CS undergrad in India.",
	"Listing out top 10 Indian dishes that everyone must try.",
	"Insights into Indian wedding rituals from recent experience.",
	"A detailed travel guide to plan a Goa trip.",
	"Hyderabadi biryani review – Is it worth the hype?",
	"Observations about how Indian parents raise their kids.",
	"My checklist for preparing for Diwali in hostel.",
	"Freelancing tips for beginners in India.",
	"Exploring whether Indian universities still teach outdated stuff.",
	"My journey exploring ancient temple architecture in Tamil Nadu.",
	"What I learned from Ramayana and Mahabharata.",
	"Tips to relieve exam stress without quitting studies.",
	"Interesting things I found while travelling in Indian trains.",
	"My experience riding cycle through Bangalore streets.",
	"How pollution affects daily life in Delhi.",
	"Top Indian influencers you should follow for productivity.",
	"Best treks in Himachal for beginners.",
	"My experience using dating apps as an Indian.",
	"Trying weird food combos with my college roommates.",
	"Planning to study abroad? Here's my advice.",
	"My personal developer journey from a college dorm in India.",
	"How yoga helped me during stressful times.",
	"Memories from our college fest – fun, chaos, and friends.",
}

var Tags = []string{
	"food", "travel", "coding", "college", "india", "festivals", "startup", "movies", "gk", "history",
	"career", "tech", "webdev", "sports", "health", "lifestyle", "parenting", "spirituality", "train", "budget",
	"culture", "education", "freelancing", "jobs", "go", "javascript", "ai", "mythology", "temple", "music",
	"biryani", "monsoon", "diwali", "holi", "love", "dating", "cs", "gate", "upsc", "bike",
	"trek", "yoga", "delhi", "kerala", "mumbai", "bangalore", "hostel", "college", "friends", "instagram",
}

var Comments = []string{
	"Awesome post!", "Loved your writing style.", "I totally relate!", "This was very helpful, thanks!",
	"Can you share more details?", "Keep posting more!", "You've inspired me to try this.", "Very informative.",
	"Thanks for the tips!", "Loved the way you explained.", "Bookmarking this!", "Will try this soon.",
	"This is gold!", "Felt nostalgic reading this.", "So well-written.", "Waiting for your next post.",
	"This deserves more views!", "Your content is underrated.", "Hats off!", "Made my day!",
	"Relatable AF!", "You spoke my mind.", "This helped me a lot.", "You should blog more.",
	"Can't wait to try it out.", "Beautifully expressed.", "So wholesome!", "Keep it up!", "Next post when?",
	"I laughed so hard!", "Thanks for sharing!", "Could you elaborate on this?", "Good one!",
	"Felt like I was there.", "Really useful.", "Great stuff!", "Learnt something new today.",
	"This is so good.", "A big thumbs up!", "Totally agree.", "Brilliant write-up.",
	"This is very accurate.", "Helpful for students.", "Enjoyed reading this.", "Perfectly put!",
	"Keep inspiring others.", "Made me emotional.", "Sharing this with friends.", "Instant follow!", "Legendary content!",
}

func Seed(db *sql.DB, userService *service.UserService, postService *service.PostService, commentService *service.CommentService, logger *zap.SugaredLogger) error {
	ctx := context.Background()

	userTemplates := generateUsers(100)
	var createdUsers []*models.User

	for _, userTemplate := range userTemplates {
		createdUser, err := userService.CreateUser(ctx, models.CreateUserRequest{
			Username: userTemplate.Username,
			Email:    userTemplate.Email,
			Password: userTemplate.PasswordHash,
		})
		if err != nil {
			logger.Error("Error creating user", zap.Error(err))
			return err
		}
		createdUsers = append(createdUsers, createdUser)
	}

	postTemplates := generatePosts(200, createdUsers)
	var createdPosts []*models.Post

	for _, postTemplate := range postTemplates {

		ctxWithUser := utils.SetUserID(ctx, postTemplate.UserID)

		createdPost, err := postService.CreatePost(ctxWithUser, models.CreatePostRequest{
			Title:   postTemplate.Title,
			Content: postTemplate.Content,
			Tags:    postTemplate.Tags,
		})
		if err != nil {
			logger.Error("Error creating post", zap.Error(err))
			return err
		}
		createdPosts = append(createdPosts, createdPost)
	}

	commentTemplates := generateComments(500, createdUsers, createdPosts)

	for _, commentTemplate := range commentTemplates {
		// Set user and post context for the comment creation
		ctxWithUser := utils.SetUserID(ctx, commentTemplate.UserID)
		ctxWithPost := context.WithValue(ctxWithUser, utils.PostIDKey, &models.Post{ID: commentTemplate.PostID})

		req := &models.CreateCommentRequest{
			Content: commentTemplate.Content,
		}

		if _, err := commentService.CreateComment(ctxWithPost, req); err != nil {
			logger.Error("Error creating comment", zap.Error(err))
			return err
		}
	}

	return nil
}

func generateUsers(num int) []*models.User {
	users := make([]*models.User, num)

	for i := range num {
		baseUsername := Usernames[i%len(Usernames)]
		users[i] = &models.User{
			Username:     fmt.Sprintf("%s_%d", baseUsername, i),
			Email:        fmt.Sprintf("%s_%d@example.com", baseUsername, i),
			PasswordHash: "password123",
		}
	}

	return users
}

func generatePosts(num int, users []*models.User) []*models.Post {
	posts := make([]*models.Post, num)

	for i := range num {
		user := users[rand.Intn(len(users))]
		posts[i] = &models.Post{
			Title:   Titles[rand.Intn(len(Titles))],
			Content: Contents[rand.Intn(len(Contents))],
			Tags:    []string{Tags[rand.Intn(len(Tags))]},
			UserID:  user.ID,
		}
	}

	return posts
}

func generateComments(num int, users []*models.User, posts []*models.Post) []*models.Comment {
	comments := make([]*models.Comment, num)

	for i := range num {
		user := users[rand.Intn(len(users))]
		post := posts[rand.Intn(len(posts))]
		comments[i] = &models.Comment{
			PostID:  post.ID,
			UserID:  user.ID,
			Content: Comments[rand.Intn(len(Comments))],
		}
	}

	return comments
}
