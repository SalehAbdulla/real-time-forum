
LOG_LEVEL value	Effective slog level
(unset/empty)	slog.LevelInfo (default)
"debug"	slog.LevelDebug
"warn"	slog.LevelWarn
"error"	slog.LevelError
How to set LOG_LEVEL
For one run:

LOG_LEVEL=debug go run ./cmd
For the current shell session:

export LOG_LEVEL=debug
go run ./cmd

Persistently (add to your .zshrc or .bash_profile):
echo 'export LOG_LEVEL=debug' >> ~/.zshrc
source ~/.zshrc