package trustedsubnet

import (
	"context"
	"log/slog"
	"net"

	"github.com/ArtShib/urlshortener/internal/lib/loghelper"
	"github.com/ArtShib/urlshortener/internal/model"
)

// TrustedSubnet структура для TrustedSubnet
type TrustedSubnet struct {
	logger     *slog.Logger
	trustedNet *net.IPNet
}

// New конструктор TrustedSubnet
func New(ctx context.Context, logger *slog.Logger, cfg *model.ConfigTrustedSubnet) *TrustedSubnet {
	log := loghelper.New(logger, "lib.trustedSubnet.New")

	trustedSubnet := &TrustedSubnet{
		logger: logger,
	}
	if cfg.Subnet == "" {
		trustedSubnet.trustedNet = nil
		log.LogError(ctx, "New", model.ErrSubnetIsEmpty)
		return trustedSubnet
	}

	_, subNet, err := net.ParseCIDR(cfg.Subnet)
	if err != nil {
		trustedSubnet.trustedNet = nil
		log.LogError(ctx, "invalid trusted subnet CIDR", err)
	} else {
		trustedSubnet.trustedNet = subNet
	}

	return trustedSubnet
}

// IsTrustedSubnet метод для определения TrustedSubnet
func (t *TrustedSubnet) IsTrustedSubnet(ctx context.Context, ipClient string) bool {
	log := loghelper.New(t.logger, "lib.trustedSubnet.IsTrustedSubnet")

	if ipClient == "" {
		log.LogDebug(ctx, "IP client is empty")
		return false
	}
	if t.trustedNet == nil {
		log.LogDebug(ctx, "subnet is empty")
		return false
	}

	ip := net.ParseIP(ipClient)
	if ip == nil {
		log.LogError(ctx, "IsTrustedSubnet", model.ErrInvalidIPFormat)
		return false
	}

	if t.trustedNet.Contains(ip) {
		return true
	}

	return false
}
