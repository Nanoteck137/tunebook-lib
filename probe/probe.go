package probe

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"gopkg.in/vansante/go-ffprobe.v2"
)

func convertMapKeysToLowercase(m map[string]any) map[string]any {
	res := make(map[string]any)
	for k, v := range m {
		res[strings.ToLower(k)] = v
	}

	return res
}

type ProbeResult struct {
	Tags     ffprobe.Tags
	Duration time.Duration
}

func ProbeMedia(ctx context.Context, filepath string) (*ProbeResult, error) {
	probe, err := ffprobe.ProbeURL(ctx, filepath)
	if err != nil {
		return nil, err
	}

	var tags ffprobe.Tags
	hasGlobalTags := probe.Format.FormatName != "ogg"

	audioStream := probe.FirstAudioStream()
	if audioStream == nil {
		// TODO(patrik): Better error?
		return nil, errors.New("contains no audio streams")
	}

	if hasGlobalTags {
		tags = probe.Format.TagList
	} else {
		tags = audioStream.TagList
	}

	tags = convertMapKeysToLowercase(tags)

	dur, err := strconv.ParseFloat(audioStream.Duration, 32)
	if err != nil {
		return nil, err
	}

	duration := time.Duration(dur * float64(time.Second))

	return &ProbeResult{
		Tags:        tags,
		Duration:    duration,
	}, nil
}
