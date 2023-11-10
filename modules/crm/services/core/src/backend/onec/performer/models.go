package performer

type performerInMongo struct {
	UUID           string `bson:"uuid,omitempty"`
	DepartmentUUID string `bson:"departmentUUID"`
	UserUUID       string `bson:"userUUID"`
}
