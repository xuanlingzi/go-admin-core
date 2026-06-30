package writer

import (
	"path/filepath"
	"testing"
	"time"
)

func TestFileWriterFilenameWithFixedNameAndRotate(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name   string
		opts   Options
		num    uint
		expect string
	}{
		{
			name: "fixed filename without rotate",
			opts: Options{
				path:     "/tmp/logs",
				filename: "app.log",
				suffix:   "log",
			},
			expect: filepath.Join("/tmp/logs", "app.log"),
		},
		{
			name: "fixed filename with day rotate",
			opts: Options{
				path:         "/tmp/logs",
				filename:     "app.log",
				suffix:       "log",
				rotatePeriod: RotatePeriodDay,
			},
			expect: filepath.Join("/tmp/logs", "app."+now.Format("2006-01-02")+".log"),
		},
		{
			name: "fixed filename with hour rotate",
			opts: Options{
				path:         "/tmp/logs",
				filename:     "app.log",
				suffix:       "log",
				rotatePeriod: RotatePeriodHour,
			},
			expect: filepath.Join("/tmp/logs", "app."+now.Format("2006-01-02-15")+".log"),
		},
		{
			name: "fixed filename with cap",
			opts: Options{
				path:     "/tmp/logs",
				filename: "app.log",
				suffix:   "log",
				cap:      1024,
			},
			num:    2,
			expect: filepath.Join("/tmp/logs", "app-[2].log"),
		},
		{
			name: "default filename with hour rotate",
			opts: Options{
				path:         "/tmp/logs",
				suffix:       "log",
				rotatePeriod: RotatePeriodHour,
			},
			expect: filepath.Join("/tmp/logs", now.Format("2006-01-02-15")+".log"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &FileWriter{
				opts: tt.opts,
				num:  tt.num,
			}
			if got := w.getFilename(); got != tt.expect {
				t.Fatalf("getFilename() = %q, want %q", got, tt.expect)
			}
		})
	}
}
