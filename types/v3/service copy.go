package types

// import (
// 	_ "embed"
// 	"strings"
// 	"text/template"

// 	option "github.com/adlerhurst/protoc-gen-go-cli/gen/proto/adlerhurst/cli/v1alpha"

// 	"google.golang.org/protobuf/compiler/protogen"
// 	"google.golang.org/protobuf/proto"
// )

// type Service struct {
// 	*protogen.Service
// 	Methods []*Method
// }

// func ServiceFromProto(service *protogen.Service) *Service {
// 	svc := &Service{
// 		Service: service,
// 		Methods: make([]*Method, 0, len(service.Methods)),
// 	}

// 	for _, method := range service.Methods {
// 		svc.Methods = append(svc.Methods, MethodFromProto(svc, method))
// 	}

// 	return svc
// }

// func (svc *Service) Generate(plugin *protogen.Plugin, file *protogen.File) (generated []*protogen.GeneratedFile, err error) {
// 	gen := plugin.NewGeneratedFile(svc.filename(file.GeneratedFilenamePrefix), file.GoImportPath)
// 	generated = append(generated, gen)

// 	header(gen, file)
// 	svc.imports(gen)
// 	if err = executeTemplate(gen, serviceTemplate, svc); err != nil {
// 		return nil, err
// 	}

// 	for _, method := range svc.Methods {
// 		methodGen, err := method.Generate(plugin, file)
// 		if err != nil {
// 			return nil, err
// 		}
// 		generated = append(generated, methodGen...)
// 	}

// 	return generated, nil
// }

// func (svc *Service) Use() string {
// 	return lower.String(svc.name())
// }

// func (svc *Service) Public() string {
// 	return title.String(svc.name())
// }

// func (svc *Service) Short() string {
// 	return string(svc.Comments.Leading)
// }

// func (svc *Service) Long() string {
// 	return string(svc.Comments.Leading) + string(svc.Comments.Trailing)
// }

// func (svc *Service) VarName() string {
// 	return svc.Public() + "Cmd"
// }
