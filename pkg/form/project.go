package form

// _______________________________________    _________.______________________ _______________.___.
// \______   \_   _____/\______   \_____  \  /   _____/|   \__    ___/\_____  \\______   \__  |   |
//  |       _/|    __)_  |     ___//   |   \ \_____  \ |   | |    |    /   |   \|       _//   |   |
//  |    |   \|        \ |    |   /    |    \/        \|   | |    |   /    |    \    |   \\____   |
//  |____|_  /_______  / |____|   \_______  /_______  /|___| |____|   \_______  /____|_  // ______|
//         \/        \/                   \/        \/                        \/       \/ \/

type CreateProject struct {
	UserID          int64  `binding:"Required"`
	Name            string `binding:"Required;AlphaDashDot;MaxSize(100)"`
	Private         bool
	Description     string `binding:"MaxSize(255)"`
	AutoInit        bool
	Gitignores      string
	License         string
	Readme          string
	CreateGitRepo   bool
	ParentProjectID int64
}

type InitializeGit struct {
	AutoInit   bool
	Gitignores string
	License    string
	Readme     string
}
