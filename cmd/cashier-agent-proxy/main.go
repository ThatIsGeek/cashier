package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/user"
	"path"
	"syscall"
	"time"

	"github.com/nsheridan/cashier/proxy"
	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/agent"
)

var (
	u, _             = user.Current()
	cfg              = pflag.String("config", path.Join(u.HomeDir, ".cashier.conf"), "Path to config file")
	ca               = pflag.String("ca", "http://localhost:10000", "CA server")
	keysize          = pflag.Int("key_size", 0, "Size of key to generate. Ignored for ed25519 keys. (default 2048 for rsa keys, 256 for ecdsa keys)")
	validity         = pflag.Duration("validity", time.Hour*24, "Key lifetime. May be overridden by the CA at signing time")
	keytype          = pflag.String("key_type", "", "Type of private key to generate - rsa, ecdsa or ed25519. (default \"rsa\")")
	publicFilePrefix = pflag.String("key_file_prefix", "", "Prefix for filename for public key and cert (optional, no default)")
	useGRPC          = pflag.Bool("use_grpc", false, "Use grpc (experimental)")
	foreground       = pflag.Bool("foreground", false, "Run cashier-agent-proxy in foreground")
	exec             = pflag.Bool("internal--exec", false, "Run the actual cashier-agent-proxy")
	curAgentSocket   = pflag.String("proxy-sock", os.Getenv("SSH_AUTH_SOCK"), "Current agent socket")
	agentSocket      = pflag.String("sock", "/tmp/proxy.sock", "Socket proxy should listen on")
)

func printEnv() {
	fmt.Printf("SSH_AUTH_SOCK=%s", *agentSocket)
}

func execPath() (string, error) {
	return os.Readlink("/proc/self/exe")
}

func doProxyAgent() {
	sock, err := net.Dial("unix", *curAgentSocket)
	if err != nil {
		log.Fatalf("Error connecting to agent: %v\n", err)
	}
	defer sock.Close()
	a := agent.NewClient(sock)
	ap := proxy.NewAgentProxy(a)

	proxySock, err := net.Listen("unix", *agentSocket)
	if err != nil {
		log.Fatalf("Error while listening: %v\n", err)
	}
	defer proxySock.Close()

	if *foreground {
		printEnv()
	}

	for {
		cs, err := proxySock.Accept()
		if err != nil {
			log.Fatalf("Error accepting connection: %v\n", err)
		}
		go agent.ServeAgent(ap, cs)
	}
}

func runProxyAgent() {
	fset := pflag.CommandLine
	flags := make([]string, fset.NFlag()+1)

	i := 0
	pflag.CommandLine.Visit(func(f *pflag.Flag) {
		flags[i] = fmt.Sprintf("--%s=%s", f.Name, f.Value.String())
		i++
	})
	flags[i] = "--internal--exec=1"

	procAttr := syscall.ProcAttr{
		Env: os.Environ(),
	}

	exe, err := execPath()
	if err != nil {
		log.Fatal("Can't find executable path")
	}

	pid, err := syscall.ForkExec(exe, flags, &procAttr)
	if err != nil {
		log.Fatalf("Error spawning child process")
	}
	log.Printf("Running with pid: %d", pid)
	printEnv()
}

func main() {
	pflag.CommandLine.MarkHidden("exec")
	pflag.Parse()
	log.SetPrefix("cashier-agent-proxy: ")
	log.SetFlags(0)

	if *foreground || *exec {
		doProxyAgent()
	} else {
		runProxyAgent()
	}
}
