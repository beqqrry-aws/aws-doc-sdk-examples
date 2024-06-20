package scenarios

import (
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/demotools"
	"github.com/awsdocs/aws-doc-sdk-examples/gov2/redshift/actions"
	"log"
)

// snippet-start:[gov2.cognito-identity-provider.Resources.complete]

// Resources keeps track of AWS resources created during an example and handles
// cleanup when the example finishes.
type Resources struct {
	userPoolId       string
	userAccessTokens []string

	redshiftActor     *actions.RedshiftActions
	redshiftDataActor *actions.RedshiftDataActions

	questioner demotools.IQuestioner
}

func (resources *Resources) init(redshiftActor *actions.RedshiftActions, redshiftDataActor *actions.RedshiftDataActions, questioner demotools.IQuestioner) {
	resources.userAccessTokens = []string{}
	resources.redshiftActor = redshiftActor
	resources.redshiftDataActor = redshiftDataActor
	resources.questioner = questioner
}

// Cleanup deletes all AWS resources created during an example.
func (resources *Resources) Cleanup() {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Something went wrong during cleanup.\n%v\n", r)
			log.Println("Use the AWS Management Console to remove any remaining resources \n" +
				"that were created for this scenario.")
		}
	}()

	//wantDelete := resources.questioner.AskBool("Do you want to remove all of the AWS resources that were created "+
	//	"during this demo (y/n)?", "y")
	//if wantDelete {
	//	for _, accessToken := range resources.userAccessTokens {
	//		err := resources.cognitoActor.DeleteUser(accessToken)
	//		if err != nil {
	//			log.Println("Couldn't delete user during cleanup.")
	//			panic(err)
	//		}
	//		log.Println("Deleted user.")
	//	}
	//	triggerList := make([]actions.TriggerInfo, len(resources.triggers))
	//	for i := 0; i < len(resources.triggers); i++ {
	//		triggerList[i] = actions.TriggerInfo{Trigger: resources.triggers[i], HandlerArn: nil}
	//	}
	//	err := resources.cognitoActor.UpdateTriggers(resources.userPoolId, triggerList...)
	//	if err != nil {
	//		log.Println("Couldn't update Cognito triggers during cleanup.")
	//		panic(err)
	//	}
	//	log.Println("Removed Cognito triggers from user pool.")
	//} else {
	//	log.Println("Be sure to remove resources when you're done with them to avoid unexpected charges!")
	//}
}

// snippet-end:[gov2.cognito-identity-provider.Resources.complete]
