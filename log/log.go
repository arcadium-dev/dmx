//  Copyright 2026 arcadium.dev <info@arcadium.dev>
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package log // import "arcadium.dev/dmx/log"

import (
	"context"
	"log/slog"
)

type (
	ctxKey struct{}
)

func FromContext(ctx context.Context) *slog.Logger {
	if ctx == nil {
		return slog.New(slog.DiscardHandler)
	}

	l, ok := ctx.Value(ctxKey{}).(*slog.Logger)
	if !ok {
		return slog.New(slog.DiscardHandler)
	}

	return l
}

func ToContext(ctx context.Context, logger *slog.Logger) context.Context {
	if ctx == nil {
		return nil
	}

	return context.WithValue(ctx, ctxKey{}, logger)
}
