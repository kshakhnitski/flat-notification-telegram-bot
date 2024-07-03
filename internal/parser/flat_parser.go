package parser

import "flat_bot/internal/model"

type FlatParser interface {
	Parse() ([]model.Flat, error)
}
