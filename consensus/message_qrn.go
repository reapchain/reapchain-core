package consensus

import (
	"errors"
	"fmt"

	"github.com/reapchain/reapchain-core/types"
)

type QrnMessage struct {
	Qrn *types.Qrn
}

// ValidateBasic performs basic validation.
func (qrnMessage *QrnMessage) ValidateBasic() error {
	return qrnMessage.Qrn.ValidateBasic()
}

// String returns a string representation.
func (qrnMessage *QrnMessage) String() string {
	return fmt.Sprintf("[Qrn %v]", qrnMessage.Qrn)
}

type HasQrnMessage struct {
	Height int64
	Index  int32
}

// ValidateBasic performs basic validation.
func (m *HasQrnMessage) ValidateBasic() error {
	if m.Height < 0 {
		return errors.New("negative Height")
	}

	if m.Index < 0 {
		return errors.New("negative Index")
	}
	return nil
}

// String returns a string representation.
func (m *HasQrnMessage) String() string {
	return fmt.Sprintf("[HasQrn VI:%v V:{%v}]", m.Index, m.Height)
}
