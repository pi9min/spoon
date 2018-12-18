package spoon

// EntityBehavior defines the interface that needs to be satisfied as an Entity.
type EntityBehavior interface {
	TableName() string
	PrimaryKey() *PrimaryKey
	Indexes() Indexes
}
