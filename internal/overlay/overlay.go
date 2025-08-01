package overlay

import (
	"fmt"
	"sync"

	"gocv.io/x/gocv"
	
	"github.com/bluenviron/mediamtx/internal/gps"
)

// Engine is the main overlay processing engine.
type Engine struct {
	config          *Config
	gpsDataProvider gps.DataProvider
	textRenderer    *TextRenderer
	enabled         bool
	mu              sync.RWMutex
}

// NewEngine creates a new overlay engine.
func NewEngine(config *Config, gpsDataProvider gps.DataProvider) (*Engine, error) {
	if !config.Enabled {
		return &Engine{
			config:  config,
			enabled: false,
		}, nil
	}

	// Initialize text renderer
	textRenderer, err := NewTextRenderer(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create text renderer: %w", err)
	}

	return &Engine{
		config:          config,
		gpsDataProvider: gpsDataProvider,
		textRenderer:    textRenderer,
		enabled:         true,
	}, nil
}

// Close closes the overlay engine and releases resources.
func (e *Engine) Close() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// No need to close GPS client as it's managed by PathManager
	e.enabled = false
	return nil
}

// IsEnabled returns whether the overlay is enabled.
func (e *Engine) IsEnabled() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.enabled
}

// ProcessFrame processes a video frame and applies overlay if enabled.
// The frame data should be in a format compatible with OpenCV (e.g., H.264 decoded frame).
func (e *Engine) ProcessFrame(frameData []byte, shipName string) ([]byte, error) {
	if !e.IsEnabled() {
		return frameData, nil
	}

	// For now, we'll skip the actual frame processing as it requires
	// decoding H.264 data to raw image data first.
	// This is a placeholder for the actual implementation.
	
	// TODO: Implement H.264 frame decoding and re-encoding
	// 1. Decode H.264 frame to raw image data
	// 2. Convert to OpenCV Mat
	// 3. Apply overlay
	// 4. Re-encode to H.264
	
	return frameData, nil
}

// ProcessMatFrame processes an OpenCV Mat frame directly (for testing or when raw frames are available).
func (e *Engine) ProcessMatFrame(frame gocv.Mat, shipName string) error {
	if !e.IsEnabled() {
		return nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()

	// Get current GPS data from provider
	gpsData := e.gpsDataProvider.GetCurrentGPS()

	// Apply overlay
	return e.textRenderer.RenderOverlay(frame, gpsData, shipName)
}

// UpdateConfig updates the engine configuration.
func (e *Engine) UpdateConfig(newConfig *Config) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// If enabling/disabling overlay
	if e.enabled != newConfig.Enabled {
		if newConfig.Enabled {
			// Need to reinitialize components
			return fmt.Errorf("enabling overlay requires engine restart")
		} else {
			// Disabling - close resources
			e.textRenderer = nil
			e.enabled = false
		}
	}

	e.config = newConfig
	return nil
}

// GetGPSData returns the current GPS data (for debugging/monitoring).
func (e *Engine) GetGPSData() *gps.Data {
	if !e.IsEnabled() || e.gpsDataProvider == nil {
		return nil
	}

	e.mu.RLock()
	defer e.mu.RUnlock()
	
	return e.gpsDataProvider.GetCurrentGPS()
}