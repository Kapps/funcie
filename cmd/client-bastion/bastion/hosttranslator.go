package bastion

import (
	"context"
	"golang.org/x/exp/slog"
	"net"
	"sync"
)

var lookupHost = net.LookupHost

// HostTranslator is an interface for translating hosts from an entered host to a resolved host.
// For example, this could translate local or unspecified IPs to the docker host when running in MacOS.
type HostTranslator interface {
	// IsHostTranslationRequired indicates whether we need to translate localhost requests to an alternative host.
	// For example, if we're running in MacOS, we need to translate localhost to host.docker.internal.
	IsHostTranslationRequired(ctx context.Context) (bool, error)
	// TranslateLocalHostToResolvedHost translates a local host (...such as localhost) to the resolved host.
	// For example, if we're running in MacOS, this could translate localhost to host.docker.internal.
	// If the host does not need to be translated, the original host is returned.
	// This also handles scenarios where the host is an IP address, such as 0.0.0.0 or 127.0.0.1.
	TranslateLocalHostToResolvedHost(ctx context.Context, host string) (string, error)
}

type dockerHostTranslator struct {
	translatedHost                  string
	isDockerHostTranslationRequired bool
	checkTranslationRequiredOnce    sync.Once
}

func NewDockerHostTranslator() HostTranslator {
	return &dockerHostTranslator{}
}

func (t *dockerHostTranslator) IsHostTranslationRequired(_ context.Context) (bool, error) {
	t.checkTranslationRequiredOnce.Do(func() {
		// MacOS uses host.docker.internal
		t.setHostIfResolves("host.docker.internal")

		// Anything below here is pretty untested...

		// Some Linux distros use the static IP 172.17.0.1
		// However it's possible that this IP is used for other things, so we don't want to assume it's always the docker host
		// For now we'll leave this out, and if we find a case where it's needed, we can add it back in
		if t.translatedHost != "" {
			slog.Info("redirecting localhost requests to resolved host.", "host", t.translatedHost)
		}
	})

	return t.isDockerHostTranslationRequired, nil
}

func (t *dockerHostTranslator) setHostIfResolves(host string) {
	if t.translatedHost != "" {
		return
	}

	_, err := lookupHost(host)
	if err == nil {
		// This is a general error, not necessarily due to no such host
		// But in practice, we can't... really handle other errors.
		// So we'll just assume that if there's an error, it's because the host doesn't exist.
		t.translatedHost = host
		t.isDockerHostTranslationRequired = true
	}
}

func (t *dockerHostTranslator) TranslateLocalHostToResolvedHost(ctx context.Context, host string) (string, error) {
	translationRequired, err := t.IsHostTranslationRequired(ctx)
	if err != nil {
		return "", err
	}

	if !translationRequired {
		return host, nil
	}

	// First, check if the host is an IP address.
	ip := net.ParseIP(host)
	if ip != nil {
		// If it is an IP, we translate if it's a loopback address (127.x.x.x) or an unspecified address (0.0.0.0).
		if ip.IsLoopback() || ip.IsUnspecified() {
			return t.translatedHost, nil
		} else {
			// Otherwise, we return the original IP address.
			return host, nil
		}
	}

	// If it's not an IP address, we translate if it's localhost.
	// This could certainly be improved to handle other cases.
	if host == "localhost" {
		return t.translatedHost, nil
	}

	// Otherwise, we return the original host.
	return host, nil
}
