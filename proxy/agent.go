package proxy

import (
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

//AgentProxy does blah blah
type AgentProxy struct {
	agent agent.Agent
}

//NewAgentProxy create new proxy
func NewAgentProxy(agent agent.Agent) AgentProxy {
	return AgentProxy{agent}
}

// List returns the identities known to the agent.
func (ap AgentProxy) List() ([]*agent.Key, error) {
	return ap.agent.List()
}

// Sign has the agent sign the data using a protocol 2 key as defined
// in [PROTOCOL.agent] section 2.6.2.
func (ap AgentProxy) Sign(key ssh.PublicKey, data []byte) (*ssh.Signature, error) {
	return ap.agent.Sign(key, data)
}

// Add adds a private key to the agent.
func (ap AgentProxy) Add(key agent.AddedKey) error {
	return ap.agent.Add(key)
}

// Remove removes all identities with the given public key.
func (ap AgentProxy) Remove(key ssh.PublicKey) error {
	return ap.agent.Remove(key)
}

// RemoveAll removes all identities.
func (ap AgentProxy) RemoveAll() error {
	return ap.agent.RemoveAll()
}

// Lock locks the agent. Sign and Remove will fail, and List will empty an empty list.
func (ap AgentProxy) Lock(passphrase []byte) error {
	return ap.agent.Lock(passphrase)
}

// Unlock undoes the effect of Lock
func (ap AgentProxy) Unlock(passphrase []byte) error {
	return ap.agent.Unlock(passphrase)
}

// Signers returns signers for all the known keys.
func (ap AgentProxy) Signers() ([]ssh.Signer, error) {
	return ap.agent.Signers()
}
