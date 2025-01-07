package teamtable

type Formatter interface {
	Format(teamTable *TeamTable) (formatted string)
}
