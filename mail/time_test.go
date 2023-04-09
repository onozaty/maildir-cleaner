package mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMailTime(t *testing.T) {

	// ARRANGE
	fileName := "1674617693.M958571P8888.localhost.localdomain,S=545,W=562:2,S"

	// ACT
	mailTime := MailTime(fileName)

	// ASSERT
	assert.Equal(t, int64(1674617693), mailTime.Unix())
}

func TestMailTime_NonDot(t *testing.T) {

	// ARRANGE
	fileName := "1674617693"

	// ACT
	mailTime := MailTime(fileName)

	// ASSERT
	assert.Equal(t, int64(1674617693), mailTime.Unix())
}

func TestMailTime_NonTime(t *testing.T) {

	// ARRANGE
	fileName := "abc"

	// ACT
	mailTime := MailTime(fileName)

	// ASSERT
	assert.Equal(t, int64(0), mailTime.Unix())
}
