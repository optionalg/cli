package v2action

type SSHInfo struct {
	SSHHostKeyFingerprint string
	SSHEndpoint           string
	SSHPasscode           string
}

func (actor Actor) GetSSHInfo() (SSHInfo, error) {
	passcode, err := actor.UAAClient.GetSSHPasscode(actor.Config)
	if err != nil {
		return SSHInfo{}, err
	}

	return SSHInfo{
		SSHEndpoint:           actor.CloudControllerClient.AppSSHEndpoint(),
		SSHHostKeyFingerprint: actor.CloudControllerClient.AppSSHHostKeyFingerprint(),
		SSHPasscode:           passcode,
	}, nil
}
