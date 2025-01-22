package stuff

import (
	"context"
	"go.opentelemetry.io/otel"
)

func DoSomeWork(ctx context.Context) {
	tracer := otel.Tracer("stuff")
	ctx1, span := tracer.Start(ctx, "DoSomeWork")
	defer span.End()
	Wait()

	tracer2 := otel.Tracer("get data in database")
	_, span2 := tracer2.Start(ctx1, "GetDataInDatabase")
	defer span2.End()

	Wait()
}
