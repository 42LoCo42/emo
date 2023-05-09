package util

// #include <mpv/client.h>
import "C"

import (
	"github.com/gen2brain/go-mpv"
)

const (
	STOP_REASON_EOF   = C.MPV_END_FILE_REASON_EOF
	STOP_REASON_STOP  = C.MPV_END_FILE_REASON_STOP
	STOP_REASON_ERROR = C.MPV_END_FILE_REASON_ERROR
)

type Mpv struct {
	*mpv.Mpv
	callbacks map[string]func(any)
	OnStop    func(reason int)
}

func (m *Mpv) Observe(
	name string,
	format mpv.Format,
	callback func(any),
) {
	if m.callbacks == nil {
		m.callbacks = map[string]func(any){}
	}

	m.ObserveProperty(0, name, format)
	m.callbacks[name] = callback
}

func (m *Mpv) Run() {
	for {
		e := m.WaitEvent(-1)

		if e.Event_Id == mpv.EVENT_END_FILE && m.OnStop != nil {
			reason := (*(*C.struct_mpv_event_end_file)(e.Data)).reason
			m.OnStop(int(reason))
		}

		name, data := readEvent(e)
		callback, ok := m.callbacks[name]
		if ok {
			callback(data)
		}
	}
}

func readEvent(e *mpv.Event) (string, any) {
	if e.Event_Id != mpv.EVENT_PROPERTY_CHANGE {
		return "", nil
	}

	property := (*C.mpv_event_property)(e.Data)
	name := C.GoString(property.name)
	format := mpv.Format(property.format)

	switch format {
	case mpv.FORMAT_FLAG:
		raw := *(*C.int)(property.data)
		if raw == 0 {
			return name, false
		} else if raw == 1 {
			return name, true
		}
	case mpv.FORMAT_INT64:
		return name, int64(*(*C.int64_t)(property.data))
	case mpv.FORMAT_DOUBLE:
		return name, float64(*(*C.double)(property.data))
	case mpv.FORMAT_STRING:
		return name, C.GoString(*(**C.char)(property.data))
	}

	return "", nil
}
