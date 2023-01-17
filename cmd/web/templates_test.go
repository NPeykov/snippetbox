package main

import (
	"testing"
	"time"

	"github.com/NPeykov/snippetbox/internal/assert"
)

func TestHumanData(t *testing.T) {
    tests := []struct {
        name string
        tm time.Time
        expect string
    }{
        {
            name : "UTC",
            tm   : time.Date(2023, 1, 17, 10, 15, 0, 0, time.UTC),
            expect : "17 Jan 2023 at 10:15",
        },
        {
            name : "Empty",
            tm : time.Time{},
            expect :  "",
        },
        {
            name : "CET",
            tm : time.Date(2023, 1, 17, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
            expect : "17 Jan 2023 at 09:15",
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            got := humanDate(test.tm)
            assert.Equal(t, got, test.expect)
        })
    }
}
