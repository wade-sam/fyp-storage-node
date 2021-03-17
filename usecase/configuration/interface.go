package configuration

//"github.com/wade-sam/fypstoragenode/entity"

type Repository interface {
	SetStorageNode(name string) error
	GetStorageNode() (string, error)
}

type Usecase interface {
	GetStorageNode() (string, error)
	SetStorageNode(name string) error
}
