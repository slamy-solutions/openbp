package runtime

import (
	"log/slog"

	"github.com/slamy-solutions/openbp/modules/runtime/libs/golang/manager/runtime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RuntimeInMongo struct {
	Namespace string `bson:"namespace"`
	Name      string `bson:"name"`
	Run       bool   `bson:"run"`

	BinaryFile primitive.ObjectID `bson:"binaryFile,omitempty"`
}

func (r *RuntimeInMongo) ToLog() slog.Attr {
	return r.ToLogWithKey("runtime")
}

func (r *RuntimeInMongo) ToLogWithKey(key string) slog.Attr {
	return slog.Group(
		key,
		slog.String("namespace", r.Namespace),
		slog.String("name", r.Name),
		slog.Bool("run", r.Run),
	)
}

func (r *RuntimeInMongo) ToGRPCRuntime() *runtime.Runtime {
	return &runtime.Runtime{
		Namespace: r.Namespace,
		Name:      r.Name,
		Run:       r.Run,
	}
}
