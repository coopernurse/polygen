package polygenlib

type JavaGenerator struct { }

func (g JavaGenerator) GenFiles(p *Package) []File {
	files := make([]File, 0)
	for i := 0; i < len(p.Structs); i++ {
		files = append(files, g.genStructClass(p, p.Structs[i]))
	}

	for i := 0; i < len(p.Interfaces); i++ {
		iface := p.Interfaces[i]
		files = append(files, g.genServiceInterface(p, iface))
		files = append(files, g.genServiceRPCServer(p, iface))
		files = append(files, g.genServiceRPCClient(p, iface))
	}

	return files
}

func JavaFilename(s string) string {
	return s + ".java"
}

func (g JavaGenerator) genStructClass(p *Package, s Struct) File {
	return File{JavaFilename(s.Name), nil}
}

func (g JavaGenerator) genServiceInterface(p *Package, i Interface) File {
	cname := i.Name
	return File{JavaFilename(cname), nil}
}

func (g JavaGenerator) genServiceRPCServer(p *Package, i Interface) File {
	cname := i.Name + "RPCServer"
	return File{JavaFilename(cname), nil}
}

func (g JavaGenerator) genServiceRPCClient(p *Package, i Interface) File {
	cname := i.Name + "RPCClient"
	return File{JavaFilename(cname), nil}
}