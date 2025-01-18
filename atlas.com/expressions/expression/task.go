package expression

import (
	"atlas-expressions/kafka/producer"
	"context"
	tenant "github.com/Chronicle20/atlas-tenant"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"time"
)

const RevertTaskName = "expression_revert_task"

type RevertTask struct {
	l        logrus.FieldLogger
	interval time.Duration
}

func NewRevertTask(l logrus.FieldLogger, interval time.Duration) *RevertTask {
	l.Infof("Initializing expression revert task to run every %dms", interval.Milliseconds())
	return &RevertTask{l, interval}
}

func (e *RevertTask) Run() {
	sctx, span := otel.GetTracerProvider().Tracer("atlas-expressions").Start(context.Background(), RevertTaskName)
	defer span.End()

	for _, exp := range GetRegistry().popExpired() {
		tctx := tenant.WithContext(sctx, exp.Tenant())
		_ = producer.ProviderImpl(e.l)(tctx)(EnvExpressionEvent)(expressionEventProvider(exp.CharacterId(), exp.WorldId(), exp.ChannelId(), exp.MapId(), 0))
	}
}

func (e *RevertTask) SleepTime() time.Duration {
	return e.interval
}
