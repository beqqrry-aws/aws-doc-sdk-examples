Python Code Example Standards

This document summarizes important points for writing and reviewing code examples written for the AWS Python (Boto3) SDK. For more information on tools and standards, see the complete list in [TCX Code Examples Standards](https://quip-amazon.com/KEBiAxE5aLO9).

1. General Structure.
    1. Service folders should include a wrapper class, test folder, and any scenarios in their own files. [Example](https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/python/example_code/keyspaces).
    2. Service folders should include a requirements.txt that includes all dependencies.
    3. Scenario files.
        1. Scenarios should include a demo_tools import for shared functionality such as handling user input.
        2. Scenarios should include an init function that takes the service wrapper as a parameter.
        3. Scenarios should include a run_scenario method that runs the scenario.
        4. Scenario files should begin with a descriptive comment block describing the steps of the scenario.
        ```
            # Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
            # SPDX-License-Identifier: Apache-2.0
            
            import logging
            import time
            import urllib.request
            import uuid
            
            import boto3
            from alive_progress import alive_bar
            from rich.console import Console
            
            from elastic_ip import ElasticIpWrapper
            from instance import EC2InstanceWrapper
            from key_pair import KeyPairWrapper
            from security_group import SecurityGroupWrapper
            
            logger = logging.getLogger(__name__)
            console = Console()
            
            
            # snippet-start:[python.example_code.ec2.Scenario_GetStartedInstances]
            class EC2InstanceScenario:
                """
                A scenario that demonstrates how to use Boto3 to manage Amazon EC2 resources.
                Covers creating a key pair, security group, launching an instance, associating
                an Elastic IP, and cleaning up resources.
                """
            
                def __init__(
                    self,
                    inst_wrapper: EC2InstanceWrapper,
                    key_wrapper: KeyPairWrapper,
                    sg_wrapper: SecurityGroupWrapper,
                    eip_wrapper: ElasticIpWrapper,
                    ssm_client: boto3.client,
                    remote_exec: bool = False,
                ):
                    """
                    Initializes the EC2InstanceScenario with the necessary AWS service wrappers.
            
                    :param inst_wrapper: Wrapper for EC2 instance operations.
                    :param key_wrapper: Wrapper for key pair operations.
                    :param sg_wrapper: Wrapper for security group operations.
                    :param eip_wrapper: Wrapper for Elastic IP operations.
                    :param ssm_client: Boto3 client for accessing SSM to retrieve AMIs.
                    :param remote_exec: Flag to indicate if the scenario is running in a remote execution
                                        environment. Defaults to False. If True, the script won't prompt
                                        for user interaction.
                    """
                    self.inst_wrapper = inst_wrapper
                    self.key_wrapper = key_wrapper
                    self.sg_wrapper = sg_wrapper
                    self.eip_wrapper = eip_wrapper
                    self.ssm_client = ssm_client
                    self.remote_exec = remote_exec
            
                def create_and_list_key_pairs(self) -> None:
                    """
                    Creates an RSA key pair for SSH access to the EC2 instance and lists available key pairs.
                    """
                    console.print("**Step 1: Create a Secure Key Pair**", style="bold cyan")
                    console.print(
                        "Let's create a secure RSA key pair for connecting to your EC2 instance."
                    )
                    key_name = f"MyUniqueKeyPair-{uuid.uuid4().hex[:8]}"
                    console.print(f"- **Key Pair Name**: {key_name}")
            
                    # Create the key pair and simulate the process with a progress bar.
                    with alive_bar(1, title="Creating Key Pair") as bar:
                        self.key_wrapper.create(key_name)
                        time.sleep(0.4)  # Simulate the delay in key creation
                        bar()
            
                    console.print(f"- **Private Key Saved to**: {self.key_wrapper.key_file_path}\n")
            
                    # List key pairs (simulated) and show a progress bar.
                    list_keys = True
                    if list_keys:
                        console.print("- Listing your key pairs...")
                        start_time = time.time()
                        with alive_bar(100, title="Listing Key Pairs") as bar:
                            while time.time() - start_time < 2:
                                time.sleep(0.2)
                                bar(10)
                            self.key_wrapper.list(5)
                            if time.time() - start_time > 2:
                                console.print(
                                    "Taking longer than expected! Please wait...",
                                    style="bold yellow",
                                )
            
                def create_security_group(self) -> None:
                    """
                    Creates a security group that controls access to the EC2 instance and adds a rule
                    to allow SSH access from the user's current public IP address.
                    """
                    console.print("**Step 2: Create a Security Group**", style="bold cyan")
                    console.print(
                        "Security groups manage access to your instance. Let's create one."
                    )
                    sg_name = f"MySecurityGroup-{uuid.uuid4().hex[:8]}"
                    console.print(f"- **Security Group Name**: {sg_name}")
            
                    # Create the security group and simulate the process with a progress bar.
                    with alive_bar(1, title="Creating Security Group") as bar:
                        self.sg_wrapper.create(
                            sg_name, "Security group for example: get started with instances."
                        )
                        time.sleep(0.5)
                        bar()
            
                    console.print(f"- **Security Group ID**: {self.sg_wrapper.security_group}\n")
            
                    # Get the current public IP to set up SSH access.
                    ip_response = urllib.request.urlopen("http://checkip.amazonaws.com")
                    current_ip_address = ip_response.read().decode("utf-8").strip()
                    console.print(
                        "Let's add a rule to allow SSH only from your current IP address."
                    )
                    console.print(f"- **Your Public IP Address**: {current_ip_address}")
                    console.print("- Automatically adding SSH rule...")
            
                    # Update security group rules to allow SSH and simulate with a progress bar.
                    with alive_bar(1, title="Updating Security Group Rules") as bar:
                        response = self.sg_wrapper.authorize_ingress(current_ip_address)
                        time.sleep(0.4)
                        if response and response.get("Return"):
                            console.print("- **Security Group Rules Updated**.")
                        else:
                            console.print(
                                "- **Error**: Couldn't update security group rules.",
                                style="bold red",
                            )
                        bar()
            
                    self.sg_wrapper.describe(self.sg_wrapper.security_group)
            
                def create_instance(self) -> None:
                    """
                    Launches an EC2 instance using an Amazon Linux 2 AMI and the created key pair
                    and security group. Displays instance details and SSH connection information.
                    """
                    # Retrieve Amazon Linux 2 AMIs from SSM.
                    ami_paginator = self.ssm_client.get_paginator("get_parameters_by_path")
                    ami_options = []
                    for page in ami_paginator.paginate(Path="/aws/service/ami-amazon-linux-latest"):
                        ami_options += page["Parameters"]
                    amzn2_images = self.inst_wrapper.get_images(
                        [opt["Value"] for opt in ami_options if "amzn2" in opt["Name"]]
                    )
                    console.print("\n**Step 3: Launch Your Instance**", style="bold cyan")
                    console.print(
                        "Let's create an instance from an Amazon Linux 2 AMI. Here are some options:"
                    )
                    image_choice = 0
                    console.print(f"- Selected AMI: {amzn2_images[image_choice]['ImageId']}\n")
            
                    # Display instance types compatible with the selected AMI
                    inst_types = self.inst_wrapper.get_instance_types(
                        amzn2_images[image_choice]["Architecture"]
                    )
                    inst_type_choice = 0
                    console.print(
                        f"- Selected instance type: {inst_types[inst_type_choice]['InstanceType']}\n"
                    )
            
                    console.print("Creating your instance and waiting for it to start...")
                    with alive_bar(1, title="Creating Instance") as bar:
                        self.inst_wrapper.create(
                            amzn2_images[image_choice]["ImageId"],
                            inst_types[inst_type_choice]["InstanceType"],
                            self.key_wrapper.key_pair["KeyName"],
                            [self.sg_wrapper.security_group],
                        )
                        time.sleep(21)
                        bar()
            
                    console.print(f"**Success! Your instance is ready:**\n", style="bold green")
                    self.inst_wrapper.display()
            
                    console.print(
                        "You can use SSH to connect to your instance. "
                        "If the connection attempt times out, you might have to manually update "
                        "the SSH ingress rule for your IP address in the AWS Management Console."
                    )
                    self._display_ssh_info()
            
                def _display_ssh_info(self) -> None:
                    """
                    Displays SSH connection information for the user to connect to the EC2 instance.
                    Handles the case where the instance does or does not have an associated public IP address.
                    """
                    if (
                        not self.eip_wrapper.elastic_ips
                        or not self.eip_wrapper.elastic_ips[0].allocation_id
                    ):
                        if self.inst_wrapper.instances:
                            instance = self.inst_wrapper.instances[0]
                            instance_id = instance["InstanceId"]
            
                            waiter = self.inst_wrapper.ec2_client.get_waiter("instance_running")
                            console.print(
                                "Waiting for the instance to be in a running state with a public IP...",
                                style="bold cyan",
                            )
            
                            with alive_bar(1, title="Waiting for Instance to Start") as bar:
                                waiter.wait(InstanceIds=[instance_id])
                                time.sleep(20)
                                bar()
            
                            instance = self.inst_wrapper.ec2_client.describe_instances(
                                InstanceIds=[instance_id]
                            )["Reservations"][0]["Instances"][0]
            
                            public_ip = instance.get("PublicIpAddress")
                            if public_ip:
                                console.print(
                                    "\nTo connect via SSH, open another command prompt and run the following command:",
                                    style="bold cyan",
                                )
                                console.print(
                                    f"\tssh -i {self.key_wrapper.key_file_path} ec2-user@{public_ip}"
                                )
                            else:
                                console.print(
                                    "Instance does not have a public IP address assigned.",
                                    style="bold red",
                                )
                        else:
                            console.print(
                                "No instance available to retrieve public IP address.",
                                style="bold red",
                            )
                    else:
                        elastic_ip = self.eip_wrapper.elastic_ips[0]
                        elastic_ip_address = elastic_ip.public_ip
                        console.print(
                            f"\tssh -i {self.key_wrapper.key_file_path} ec2-user@{elastic_ip_address}"
                        )
            
                    if not self.remote_exec:
                        console.print("\nOpen a new terminal tab to try the above SSH command.")
                        input("Press Enter to continue...")
            
                def associate_elastic_ip(self) -> None:
                    """
                    Allocates an Elastic IP address and associates it with the EC2 instance.
                    Displays the Elastic IP address and SSH connection information.
                    """
                    console.print("\n**Step 4: Allocate an Elastic IP Address**", style="bold cyan")
                    console.print(
                        "You can allocate an Elastic IP address and associate it with your instance\n"
                        "to keep a consistent IP address even when your instance restarts."
                    )
            
                    with alive_bar(1, title="Allocating Elastic IP") as bar:
                        elastic_ip = self.eip_wrapper.allocate()
                        time.sleep(0.5)
                        bar()
            
                    console.print(
                        f"- **Allocated Static Elastic IP Address**: {elastic_ip.public_ip}."
                    )
            
                    with alive_bar(1, title="Associating Elastic IP") as bar:
                        self.eip_wrapper.associate(
                            elastic_ip.allocation_id, self.inst_wrapper.instances[0]["InstanceId"]
                        )
                        time.sleep(2)
                        bar()
            
                    console.print(f"- **Associated Elastic IP with Your Instance**.")
                    console.print(
                        "You can now use SSH to connect to your instance by using the Elastic IP."
                    )
                    self._display_ssh_info()
            
                def stop_and_start_instance(self) -> None:
                    """
                    Stops and restarts the EC2 instance. Displays instance state and explains
                    changes that occur when the instance is restarted, such as the potential change
                    in the public IP address unless an Elastic IP is associated.
                    """
                    console.print("\n**Step 5: Stop and Start Your Instance**", style="bold cyan")
                    console.print("Let's stop and start your instance to see what changes.")
                    console.print("- **Stopping your instance and waiting until it's stopped...**")
            
                    with alive_bar(1, title="Stopping Instance") as bar:
                        self.inst_wrapper.stop()
                        time.sleep(360)
                        bar()
            
                    console.print("- **Your instance is stopped. Restarting...**")
            
                    with alive_bar(1, title="Starting Instance") as bar:
                        self.inst_wrapper.start()
                        time.sleep(20)
                        bar()
            
                    console.print("**Your instance is running.**", style="bold green")
                    self.inst_wrapper.display()
            
                    elastic_ip = (
                        self.eip_wrapper.elastic_ips[0] if self.eip_wrapper.elastic_ips else None
                    )
            
                    if elastic_ip is None or elastic_ip.allocation_id is None:
                        console.print(
                            "- **Note**: Every time your instance is restarted, its public IP address changes."
                        )
                    else:
                        console.print(
                            f"Because you have associated an Elastic IP with your instance, you can \n"
                            f"connect by using a consistent IP address after the instance restarts: {elastic_ip.public_ip}"
                        )
            
                    self._display_ssh_info()
            
                def cleanup(self) -> None:
                    """
                    Cleans up all the resources created during the scenario, including disassociating
                    and releasing the Elastic IP, terminating the instance, deleting the security
                    group, and deleting the key pair.
                    """
                    console.print("\n**Step 6: Clean Up Resources**", style="bold cyan")
                    console.print("Cleaning up resources:")
            
                    for elastic_ip in self.eip_wrapper.elastic_ips:
                        console.print(f"- **Elastic IP**: {elastic_ip.public_ip}")
            
                        with alive_bar(1, title="Disassociating Elastic IP") as bar:
                            self.eip_wrapper.disassociate(elastic_ip.allocation_id)
                            time.sleep(2)
                            bar()
            
                        console.print("\t- **Disassociated Elastic IP from the Instance**")
            
                        with alive_bar(1, title="Releasing Elastic IP") as bar:
                            self.eip_wrapper.release(elastic_ip.allocation_id)
                            time.sleep(1)
                            bar()
            
                        console.print("\t- **Released Elastic IP**")
            
                    console.print(f"- **Instance**: {self.inst_wrapper.instances[0]['InstanceId']}")
            
                    with alive_bar(1, title="Terminating Instance") as bar:
                        self.inst_wrapper.terminate()
                        time.sleep(380)
                        bar()
            
                    console.print("\t- **Terminated Instance**")
            
                    console.print(f"- **Security Group**: {self.sg_wrapper.security_group}")
            
                    with alive_bar(1, title="Deleting Security Group") as bar:
                        self.sg_wrapper.delete(self.sg_wrapper.security_group)
                        time.sleep(1)
                        bar()
            
                    console.print("\t- **Deleted Security Group**")
            
                    console.print(f"- **Key Pair**: {self.key_wrapper.key_pair['KeyName']}")
            
                    with alive_bar(1, title="Deleting Key Pair") as bar:
                        self.key_wrapper.delete(self.key_wrapper.key_pair["KeyName"])
                        time.sleep(0.4)
                        bar()
            
                    console.print("\t- **Deleted Key Pair**")
            
                def run_scenario(self) -> None:
                    """
                    Executes the entire EC2 instance scenario: creates key pairs, security groups,
                    launches an instance, associates an Elastic IP, and cleans up all resources.
                    """
                    logging.basicConfig(level=logging.INFO, format="%(levelname)s: %(message)s")
            
                    console.print("-" * 88)
                    console.print(
                        "Welcome to the Amazon Elastic Compute Cloud (Amazon EC2) get started with instances demo.",
                        style="bold magenta",
                    )
                    console.print("-" * 88)
            
                    self.create_and_list_key_pairs()
                    self.create_security_group()
                    self.create_instance()
                    self.stop_and_start_instance()
                    self.associate_elastic_ip()
                    self.stop_and_start_instance()
                    self.cleanup()
            
                    console.print("\nThanks for watching!", style="bold green")
                    console.print("-" * 88)
            
            
            if __name__ == "__main__":
                try:
                    scenario = EC2InstanceScenario(
                        EC2InstanceWrapper.from_client(),
                        KeyPairWrapper.from_client(),
                        SecurityGroupWrapper.from_client(),
                        ElasticIpWrapper.from_client(),
                        boto3.client("ssm"),
                    )
                    scenario.run_scenario()
                except Exception:
                    logging.exception("Something went wrong with the demo.")
            # snippet-end:[python.example_code.ec2.Scenario_GetStartedInstances]```

    4. Wrapper classes.
        1. Wrapper methods should provide additional context when calling service actions. For example, specify certain parameters and return true or false. Do not use Request/Response classes directly as the parameter and response types.
        2. Wrapper classes should include a from_client method that initializes the service client. This client should be reused throughout the wrapper class.
        3. Wrapper classes should include a class declaration snippet to include in the Code Library metadata.

```
            # Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
            # SPDX-License-Identifier: Apache-2.0
            import logging
            from typing import Any, Dict, List, Optional, Union
            
            import boto3
            from botocore.exceptions import ClientError
            
            logger = logging.getLogger(__name__)
            
            
            # snippet-start:[python.example_code.ec2.ElasticIpWrapper.class]
            # snippet-start:[python.example_code.ec2.ElasticIpWrapper.decl]
            class ElasticIpWrapper:
                """Encapsulates Amazon Elastic Compute Cloud (Amazon EC2) Elastic IP address actions using the client interface."""
            
                class ElasticIp:
                    """Represents an Elastic IP and its associated instance."""
            
                    def __init__(
                        self, allocation_id: str, public_ip: str, instance_id: Optional[str] = None
                    ) -> None:
                        """
                        Initializes the ElasticIp object.
            
                        :param allocation_id: The allocation ID of the Elastic IP.
                        :param public_ip: The public IP address of the Elastic IP.
                        :param instance_id: The ID of the associated EC2 instance, if any.
                        """
                        self.allocation_id = allocation_id
                        self.public_ip = public_ip
                        self.instance_id = instance_id
            
                def __init__(self, ec2_client: Any) -> None:
                    """
                    Initializes the ElasticIpWrapper with an EC2 client.
            
                    :param ec2_client: A Boto3 Amazon EC2 client. This client provides low-level
                                       access to AWS EC2 services.
                    """
                    self.ec2_client = ec2_client
                    self.elastic_ips: List[ElasticIpWrapper.ElasticIp] = []
            
                @classmethod
                def from_client(cls) -> "ElasticIpWrapper":
                    """
                    Creates an ElasticIpWrapper instance with a default EC2 client.
            
                    :return: An instance of ElasticIpWrapper initialized with the default EC2 client.
                    """
                    ec2_client = boto3.client("ec2")
                    return cls(ec2_client)
            
                # snippet-end:[python.example_code.ec2.ElasticIpWrapper.decl]
            
                # snippet-start:[python.example_code.ec2.AllocateAddress]
                def allocate(self) -> "ElasticIpWrapper.ElasticIp":
                    """
                    Allocates an Elastic IP address that can be associated with an Amazon EC2
                    instance. By using an Elastic IP address, you can keep the public IP address
                    constant even when you restart the associated instance.
            
                    :return: The ElasticIp object for the newly created Elastic IP address.
                    :raises ClientError: If the allocation fails, such as reaching the maximum limit of Elastic IPs.
                    """
                    try:
                        response = self.ec2_client.allocate_address(Domain="vpc")
                        elastic_ip = self.ElasticIp(
                            allocation_id=response["AllocationId"], public_ip=response["PublicIp"]
                        )
                        self.elastic_ips.append(elastic_ip)
                    except ClientError as err:
                        if err.response["Error"]["Code"] == "AddressLimitExceeded":
                            logger.error(
                                "Max IP's reached. Release unused addresses or contact AWS Support for an increase."
                            )
                        raise err
                    return elastic_ip
            
                # snippet-end:[python.example_code.ec2.AllocateAddress]
            
                # snippet-start:[python.example_code.ec2.AssociateAddress]
                def associate(
                    self, allocation_id: str, instance_id: str
                ) -> Union[Dict[str, Any], None]:
                    """
                    Associates an Elastic IP address with an instance. When this association is
                    created, the Elastic IP's public IP address is immediately used as the public
                    IP address of the associated instance.
            
                    :param allocation_id: The allocation ID of the Elastic IP.
                    :param instance_id: The ID of the Amazon EC2 instance.
                    :return: A response that contains the ID of the association, or None if no Elastic IP is found.
                    :raises ClientError: If the association fails, such as when the instance ID is not found.
                    """
                    elastic_ip = self.get_elastic_ip_by_allocation(self.elastic_ips, allocation_id)
                    if elastic_ip is None:
                        logger.info(f"No Elastic IP found with allocation ID {allocation_id}.")
                        return None
            
                    try:
                        response = self.ec2_client.associate_address(
                            AllocationId=allocation_id, InstanceId=instance_id
                        )
                        elastic_ip.instance_id = (
                            instance_id  # Track the instance associated with this Elastic IP.
                        )
                    except ClientError as err:
                        if err.response["Error"]["Code"] == "InvalidInstanceID.NotFound":
                            logger.error(
                                f"Failed to associate Elastic IP {allocation_id} with {instance_id} "
                                "because the specified instance ID does not exist or has not propagated fully. "
                                "Verify the instance ID and try again, or wait a few moments before attempting to "
                                "associate the Elastic IP address."
                            )
                        raise
                    return response
            
                # snippet-end:[python.example_code.ec2.AssociateAddress]
            
                # snippet-start:[python.example_code.ec2.DisassociateAddress]
                def disassociate(self, allocation_id: str) -> None:
                    """
                    Removes an association between an Elastic IP address and an instance. When the
                    association is removed, the instance is assigned a new public IP address.
            
                    :param allocation_id: The allocation ID of the Elastic IP to disassociate.
                    :raises ClientError: If the disassociation fails, such as when the association ID is not found.
                    """
                    elastic_ip = self.get_elastic_ip_by_allocation(self.elastic_ips, allocation_id)
                    if elastic_ip is None or elastic_ip.instance_id is None:
                        logger.info(
                            f"No association found for Elastic IP with allocation ID {allocation_id}."
                        )
                        return
            
                    try:
                        # Retrieve the association ID before disassociating
                        response = self.ec2_client.describe_addresses(AllocationIds=[allocation_id])
                        association_id = response["Addresses"][0].get("AssociationId")
            
                        if association_id:
                            self.ec2_client.disassociate_address(AssociationId=association_id)
                            elastic_ip.instance_id = None  # Remove the instance association
                        else:
                            logger.info(
                                f"No Association ID found for Elastic IP with allocation ID {allocation_id}."
                            )
            
                    except ClientError as err:
                        if err.response["Error"]["Code"] == "InvalidAssociationID.NotFound":
                            logger.error(
                                f"Failed to disassociate Elastic IP {allocation_id} "
                                "because the specified association ID for the Elastic IP address was not found. "
                                "Verify the association ID and ensure the Elastic IP is currently associated with a "
                                "resource before attempting to disassociate it."
                            )
                        raise
            
                # snippet-end:[python.example_code.ec2.DisassociateAddress]
            
                # snippet-start:[python.example_code.ec2.ReleaseAddress]
                def release(self, allocation_id: str) -> None:
                    """
                    Releases an Elastic IP address. After the Elastic IP address is released,
                    it can no longer be used.
            
                    :param allocation_id: The allocation ID of the Elastic IP to release.
                    :raises ClientError: If the release fails, such as when the Elastic IP address is not found.
                    """
                    elastic_ip = self.get_elastic_ip_by_allocation(self.elastic_ips, allocation_id)
                    if elastic_ip is None:
                        logger.info(f"No Elastic IP found with allocation ID {allocation_id}.")
                        return
            
                    try:
                        self.ec2_client.release_address(AllocationId=allocation_id)
                        self.elastic_ips.remove(elastic_ip)  # Remove the Elastic IP from the list
                    except ClientError as err:
                        if err.response["Error"]["Code"] == "InvalidAddress.NotFound":
                            logger.error(
                                f"Failed to release Elastic IP address {allocation_id} "
                                "because it could not be found. Verify the Elastic IP address "
                                "and ensure it is allocated to your account in the correct region "
                                "before attempting to release it."
                            )
                        raise
            
                # snippet-end:[python.example_code.ec2.ReleaseAddress]
            
                @staticmethod
                def get_elastic_ip_by_allocation(
                    elastic_ips: List["ElasticIpWrapper.ElasticIp"], allocation_id: str
                ) -> Optional["ElasticIpWrapper.ElasticIp"]:
                    """
                    Retrieves an Elastic IP object by its allocation ID from a given list of Elastic IPs.
            
                    :param elastic_ips: A list of ElasticIp objects.
                    :param allocation_id: The allocation ID of the Elastic IP to retrieve.
                    :return: The ElasticIp object associated with the allocation ID, or None if not found.
                    """
                    return next(
                        (ip for ip in elastic_ips if ip.allocation_id == allocation_id), None
                    )
            
            # snippet-end:[python.example_code.ec2.ElasticIpWrapper.class]
```

2. Language Features
    1. Comments.
        1. All classes and methods should include a descriptive comment block.
        2. All methods should include param tag comments for each parameter and the method response.
    2. Logging
        1. Use the logging module for logging rather than print statements. 
             1. Note: User interactions may still use print statements.
        2. Set up a logger at the module level and use it within classes and functions.
    3. Error Handling.
        1. Use try and except blocks to handle potential errors gracefully.
        2. Log error messages using the logger.
        3. Raise exceptions when appropriate after logging the error.
    4. Function signatures.
        1. Function signatures should include type hints.
    4. Testing.
        1. All wrapper methods and scenario methods should have test coverage.
        2. Unit tests should use the established service method stubbing patterns for testing. Here is an example:

```# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0

from botocore.exceptions import ClientError
import pytest


class MockManager:
    def __init__(self, stub_runner, scenario_data, input_mocker):
        self.scenario_data = scenario_data
        self.ks_exists = False
        self.ks_name = "test-ks"
        self.ks_arn = "arn:aws:cassandra:test-region:111122223333:/keyspace/test-ks"
        self.keyspaces = [
            {"keyspaceName": f"ks-{ind}", "resourceArn": self.ks_arn}
            for ind in range(1, 4)
        ]
        answers = [self.ks_name]
        input_mocker.mock_answers(answers)
        self.stub_runner = stub_runner

    def setup_stubs(self, error, stop_on, stubber):
        with self.stub_runner(error, stop_on) as runner:
            if self.ks_exists:
                runner.add(stubber.stub_get_keyspace, self.ks_name, self.ks_arn)
            else:
                runner.add(
                    stubber.stub_get_keyspace,
                    self.ks_name,
                    self.ks_arn,
                    error_code="ResourceNotFoundException",
                )
                runner.add(stubber.stub_create_keyspace, self.ks_name, self.ks_arn)
                runner.add(stubber.stub_get_keyspace, self.ks_name, self.ks_arn)
            runner.add(stubber.stub_list_keyspaces, self.keyspaces)


@pytest.fixture
def mock_mgr(stub_runner, scenario_data, input_mocker):
    return MockManager(stub_runner, scenario_data, input_mocker)


@pytest.mark.parametrize("ks_exists", [True, False])
def test_create_keyspace(mock_mgr, capsys, ks_exists):
    mock_mgr.ks_exists = ks_exists
    mock_mgr.setup_stubs(None, None, mock_mgr.scenario_data.stubber)

    mock_mgr.scenario_data.scenario.create_keyspace()

    capt = capsys.readouterr()
    assert mock_mgr.ks_name in capt.out
    for ks in mock_mgr.keyspaces:
        assert ks["keyspaceName"] in capt.out


@pytest.mark.parametrize(
    "error, stop_on_index",
    [
        ("TESTERROR-stub_get_keyspace", 0),
        ("TESTERROR-stub_create_keyspace", 1),
        ("TESTERROR-stub_list_keyspaces", 2),
    ],
)
def test_create_keyspace_error(mock_mgr, caplog, error, stop_on_index):
    mock_mgr.setup_stubs(error, stop_on_index, mock_mgr.scenario_data.stubber)

    with pytest.raises(ClientError) as exc_info:
        mock_mgr.scenario_data.scenario.create_keyspace()
    assert exc_info.value.response["Error"]["Code"] == error

    assert error in caplog.text
```

1. Integration tests should be marked with @pytest.mark.integ.  Here is an example:


```
# Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
# SPDX-License-Identifier: Apache-2.0
import os
import pytest
import sys

import boto3

sys.path.insert(0, os.path.abspath(os.path.join(os.path.dirname(__file__), "..")))
from elastic_ip import ElasticIpWrapper
from instance import EC2InstanceWrapper
from key_pair import KeyPairWrapper
from scenario_get_started_instances import EC2InstanceScenario
from security_group import SecurityGroupWrapper


@pytest.mark.integ
def test_scenario_happy_path(capsys):
    # Instantiate the wrappers
    inst_wrapper = EC2InstanceWrapper.from_client()
    key_wrapper = KeyPairWrapper.from_client()
    sg_wrapper = SecurityGroupWrapper.from_client()
    eip_wrapper = ElasticIpWrapper.from_client()

    # Create the scenario object
    scenario = EC2InstanceScenario(
        inst_wrapper=inst_wrapper,
        key_wrapper=key_wrapper,
        sg_wrapper=sg_wrapper,
        eip_wrapper=eip_wrapper,
        ssm_client=boto3.client("ssm"),
        remote_exec=True,
    )

    # Run the scenario, exit 1 with error
    scenario.run_scenario()

    capt = capsys.readouterr()
    assert "Thanks for watching!" in capt.out
```


1. Metadata Snippet Contents
    1. Metadata for Action examples should contain at minimum the following snippets.
        1. A declaration snippet, to show the service and class setup.
        2. A snippet to show the action itself within context.
        3. If more than one variation of the Action is included, use descriptions in the metadata to explain the differences.
    2. Metadata for Scenario examples can contain the entire wrapper and scenario code, but should include descriptions for both.

```
keyspaces_CreateKeyspace:
          languages:
            Python:
              versions:
                - sdk_version: 3
                  github: python/example_code/keyspaces
                  sdkguide:
                  excerpts:
                    - description:
                      snippet_tags:
                        - python.example_code.keyspaces.KeyspaceWrapper.decl
                        - python.example_code.keyspaces.CreateKeyspace
          services:
            keyspaces: {CreateKeyspace}
```