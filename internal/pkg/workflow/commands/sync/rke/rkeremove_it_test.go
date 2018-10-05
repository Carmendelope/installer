/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

/*
RUN_INTEGRATION_TEST=true
RKE_TARGET_NODES=172.28.128.3
RKE_PRIVATE_KEY=/Users/daniel/.ssh/id_rsa
RKE_BINARY=/Users/daniel/Downloads/rke_darwin-amd64
*/

package rke

/*
func removeVagrant(t *testing.T) {

	var rkeBinaryPath = os.Getenv("RKE_BINARY")
	var privateKeyPath = os.Getenv("RKE_PRIVATE_KEY")
	var target = os.Getenv("RKE_TARGET_NODES")
	targetNodes := []string{target}

	cmd := NewRKERemove(
		rkeBinaryPath,
		*NewClusterConfig(
			"testClusterIT",
			targetNodes[0],
			targetNodes,
			targetNodes,
			"vagrant",
			privateKeyPath), TestTemplate)

	commandHandler := handler.GetCommandHandler()
	helper := NewHandlerHelper()
	commandHandler.AddCommand(cmd.ID(), helper.resultCallback, helper.logCallback)

	result, err := cmd.Run("workflowID")
	assert.Nil(t, err, "expecting no error")
	assert.NotNil(t, result, "expecting result")
	log.Debug().Msg(result.String())
	assert.True(t, result.Success, "expecting command to succedded")
}

func TestRKERemove(t *testing.T) {
	if RunRKEInstallTest() {
		utils.EnableDebug()
		removeVagrant(t)
	} else {
		log.Info().Msg("skipping RKE remove test")
	}
}

*/