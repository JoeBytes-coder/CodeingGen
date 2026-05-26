package generators

import "configgen/internal/domain"

// Generator 定义配置生成器的接口
type Generator interface {
	Generate(req domain.ConfigRequest) (string, error)
}
