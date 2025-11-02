package alert

import (
	"fmt"
	"time"

	beep "github.com/faiface/beep"
	"github.com/faiface/beep/generators"
	"github.com/faiface/beep/speaker"
	"github.com/gen2brain/beeep"
	"github.com/raulbattistini/private-empty/internal/util/logger"
	"github.com/raulbattistini/private-empty/pkg/metrics"
)

var GlobalSampleRate = beep.SampleRate(44100)

func initSpeaker() {
	log := logger.Log()

	audioOnce.Do(func() {
		if err := speaker.Init(GlobalSampleRate, GlobalSampleRate.N(time.Millisecond*500)); err != nil {
			log.Error("Failed to initialize global audio speaker. Sound alerts disabled: %v", err)
			return
		}

		log.Info("Audio speaker initialized successfully.")
	})
}

// this version had to be added: emulating the classy Windows sound, specially with goto's (initially) and labels was too much to be avoiding
func AnnoyingDesktopNotifier(m metrics.Metric) {
	log := logger.Log()
	title := fmt.Sprintf("⚠️ Alert: %s", m.Source)
	message := fmt.Sprintf("Value: %.2f\nMessage: %s", m.Value, m.Message)

	initSpeaker()

	// Try to play sound
	if err := playAlertSound(); err != nil {
		log.Error("Failed to play sound: %v", err)
		// Continue anyway - notification is more important
	}

	// Send notification
	if err := beeep.Alert(title, message, ""); err != nil {
		log.Error("Desktop notification failed: %v", err)
	} else {
		log.Info("Sent desktop notification with sound: %s", title)
	}
}

func playAlertSound() error {
	tone, err := generators.SinTone(GlobalSampleRate, 600*beeep.DefaultFreq)
	if err != nil {
		return fmt.Errorf("create tone: %w", err)
	}

	limitedTone := beep.Take(GlobalSampleRate.N(200*time.Millisecond), tone)
	speaker.Play(limitedTone)
	return nil
}
func DesktopNotifier(m metrics.Metric) {
	log := logger.Log()

	title := fmt.Sprintf("⚠️ Alert: %s", m.Source)
	message := fmt.Sprintf("Value: %.2f\nMessage: %s", m.Value, m.Message)

	err := beeep.Alert(title, message, "")

	if err != nil {
		log.Error("Desktop notification failed: %v", err)
		fmt.Printf("[ALERT FALLBACK] %s - %s\n", title, message)
	} else {
		log.Info("Sent desktop notification: %s", title)
	}
}
