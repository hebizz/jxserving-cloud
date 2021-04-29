package interfaces

type Dataset struct {
  Name             string
  CreatedTimestamp int64
  Sets             []string
  Types            string
  SavedPath        string
}
