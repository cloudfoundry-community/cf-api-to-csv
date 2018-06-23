package main

import "time"

type appsAPIResponse struct {
	TotalResults int         `json:"total_results"`
	TotalPages   int         `json:"total_pages"`
	PrevURL      interface{} `json:"prev_url"`
	NextURL      interface{} `json:"next_url"`
	Resources    []struct {
		Metadata struct {
			GUID      string    `json:"guid"`
			URL       string    `json:"url"`
			CreatedAt time.Time `json:"created_at"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"metadata"`
		Entity struct {
			Name                     string      `json:"name"`
			Production               bool        `json:"production"`
			SpaceGUID                string      `json:"space_guid"`
			StackGUID                string      `json:"stack_guid"`
			Buildpack                interface{} `json:"buildpack"`
			DetectedBuildpack        interface{} `json:"detected_buildpack"`
			DetectedBuildpackGUID    interface{} `json:"detected_buildpack_guid"`
			EnvironmentJSON          interface{} `json:"environment_json"`
			Memory                   int         `json:"memory"`
			Instances                int         `json:"instances"`
			DiskQuota                int         `json:"disk_quota"`
			State                    string      `json:"state"`
			Version                  string      `json:"version"`
			Command                  interface{} `json:"command"`
			Console                  bool        `json:"console"`
			Debug                    interface{} `json:"debug"`
			StagingTaskID            interface{} `json:"staging_task_id"`
			PackageState             string      `json:"package_state"`
			HealthCheckHTTPEndpoint  string      `json:"health_check_http_endpoint"`
			HealthCheckType          string      `json:"health_check_type"`
			HealthCheckTimeout       interface{} `json:"health_check_timeout"`
			StagingFailedReason      interface{} `json:"staging_failed_reason"`
			StagingFailedDescription interface{} `json:"staging_failed_description"`
			Diego                    bool        `json:"diego"`
			DockerImage              interface{} `json:"docker_image"`
			DockerCredentials        struct {
				Username interface{} `json:"username"`
				Password interface{} `json:"password"`
			} `json:"docker_credentials"`
			PackageUpdatedAt     time.Time   `json:"package_updated_at"`
			DetectedStartCommand string      `json:"detected_start_command"`
			EnableSSH            bool        `json:"enable_ssh"`
			Ports                interface{} `json:"ports"`
			SpaceURL             string      `json:"space_url"`
			StackURL             string      `json:"stack_url"`
			RoutesURL            string      `json:"routes_url"`
			EventsURL            string      `json:"events_url"`
			ServiceBindingsURL   string      `json:"service_bindings_url"`
			RouteMappingsURL     string      `json:"route_mappings_url"`
		} `json:"entity"`
	} `json:"resources"`
}

type appsAPIResource struct {
	Metadata struct {
		GUID      string    `json:"guid"`
		URL       string    `json:"url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Name                     string      `json:"name"`
		Production               bool        `json:"production"`
		SpaceGUID                string      `json:"space_guid"`
		StackGUID                string      `json:"stack_guid"`
		Buildpack                interface{} `json:"buildpack"`
		DetectedBuildpack        interface{} `json:"detected_buildpack"`
		DetectedBuildpackGUID    interface{} `json:"detected_buildpack_guid"`
		EnvironmentJSON          interface{} `json:"environment_json"`
		Memory                   int         `json:"memory"`
		Instances                int         `json:"instances"`
		DiskQuota                int         `json:"disk_quota"`
		State                    string      `json:"state"`
		Version                  string      `json:"version"`
		Command                  interface{} `json:"command"`
		Console                  bool        `json:"console"`
		Debug                    interface{} `json:"debug"`
		StagingTaskID            interface{} `json:"staging_task_id"`
		PackageState             string      `json:"package_state"`
		HealthCheckHTTPEndpoint  string      `json:"health_check_http_endpoint"`
		HealthCheckType          string      `json:"health_check_type"`
		HealthCheckTimeout       interface{} `json:"health_check_timeout"`
		StagingFailedReason      interface{} `json:"staging_failed_reason"`
		StagingFailedDescription interface{} `json:"staging_failed_description"`
		Diego                    bool        `json:"diego"`
		DockerImage              interface{} `json:"docker_image"`
		DockerCredentials        struct {
			Username interface{} `json:"username"`
			Password interface{} `json:"password"`
		} `json:"docker_credentials"`
		PackageUpdatedAt     time.Time   `json:"package_updated_at"`
		DetectedStartCommand string      `json:"detected_start_command"`
		EnableSSH            bool        `json:"enable_ssh"`
		Ports                interface{} `json:"ports"`
		SpaceURL             string      `json:"space_url"`
		StackURL             string      `json:"stack_url"`
		RoutesURL            string      `json:"routes_url"`
		EventsURL            string      `json:"events_url"`
		ServiceBindingsURL   string      `json:"service_bindings_url"`
		RouteMappingsURL     string      `json:"route_mappings_url"`
	} `json:"entity"`
}

type eventsAPIResponse struct {
	TotalResults int         `json:"total_results"`
	TotalPages   int         `json:"total_pages"`
	PrevURL      interface{} `json:"prev_url"`
	NextURL      interface{} `json:"next_url"`
	Resources    []struct {
		Metadata struct {
			GUID      string      `json:"guid"`
			URL       string      `json:"url"`
			CreatedAt time.Time   `json:"created_at"`
			UpdatedAt interface{} `json:"updated_at"`
		} `json:"metadata"`
		Entity struct {
			Type      string    `json:"type"`
			Actor     string    `json:"actor"`
			ActorType string    `json:"actor_type"`
			ActorName string    `json:"actor_name"`
			Actee     string    `json:"actee"`
			ActeeType string    `json:"actee_type"`
			ActeeName string    `json:"actee_name"`
			Timestamp time.Time `json:"timestamp"`
			Metadata  struct {
				Request struct {
					Name                  string `json:"name"`
					Instances             int    `json:"instances"`
					Memory                int    `json:"memory"`
					State                 string `json:"state"`
					EnvironmentJSON       string `json:"environment_json"`
					DockerCredentialsJSON string `json:"docker_credentials_json"`
				} `json:"request"`
			} `json:"metadata"`
			SpaceGUID        string `json:"space_guid"`
			OrganizationGUID string `json:"organization_guid"`
		} `json:"entity"`
	} `json:"resources"`
}

type eventsAPIResource struct {
	Metadata struct {
		GUID      string      `json:"guid"`
		URL       string      `json:"url"`
		CreatedAt time.Time   `json:"created_at"`
		UpdatedAt interface{} `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Type      string    `json:"type"`
		Actor     string    `json:"actor"`
		ActorType string    `json:"actor_type"`
		ActorName string    `json:"actor_name"`
		Actee     string    `json:"actee"`
		ActeeType string    `json:"actee_type"`
		ActeeName string    `json:"actee_name"`
		Timestamp time.Time `json:"timestamp"`
		Metadata  struct {
			Request struct {
				Name                  string `json:"name"`
				Instances             int    `json:"instances"`
				Memory                int    `json:"memory"`
				State                 string `json:"state"`
				EnvironmentJSON       string `json:"environment_json"`
				DockerCredentialsJSON string `json:"docker_credentials_json"`
			} `json:"request"`
		} `json:"metadata"`
		SpaceGUID        string `json:"space_guid"`
		OrganizationGUID string `json:"organization_guid"`
	} `json:"entity"`
}

type serviceBindingsAPIResponse struct {
	TotalResults int         `json:"total_results"`
	TotalPages   int         `json:"total_pages"`
	PrevURL      interface{} `json:"prev_url"`
	NextURL      interface{} `json:"next_url"`
	Resources    []struct {
		Metadata struct {
			GUID      string      `json:"guid"`
			URL       string      `json:"url"`
			CreatedAt time.Time   `json:"created_at"`
			UpdatedAt interface{} `json:"updated_at"`
		} `json:"metadata"`
		Entity struct {
			AppGUID             string `json:"app_guid"`
			ServiceInstanceGUID string `json:"service_instance_guid"`
			Credentials         struct {
				CredsKey100 string `json:"creds-key-100"`
			} `json:"credentials"`
			BindingOptions struct {
			} `json:"binding_options"`
			GatewayData        interface{} `json:"gateway_data"`
			GatewayName        string      `json:"gateway_name"`
			SyslogDrainURL     interface{} `json:"syslog_drain_url"`
			AppURL             string      `json:"app_url"`
			ServiceInstanceURL string      `json:"service_instance_url"`
		} `json:"entity"`
	} `json:"resources"`
}

type serviceBindingsAPIResource struct {
	Metadata struct {
		GUID      string      `json:"guid"`
		URL       string      `json:"url"`
		CreatedAt time.Time   `json:"created_at"`
		UpdatedAt interface{} `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		AppGUID             string `json:"app_guid"`
		ServiceInstanceGUID string `json:"service_instance_guid"`
		Credentials         struct {
			CredsKey100 string `json:"creds-key-100"`
		} `json:"credentials"`
		BindingOptions struct {
		} `json:"binding_options"`
		GatewayData        interface{} `json:"gateway_data"`
		GatewayName        string      `json:"gateway_name"`
		SyslogDrainURL     interface{} `json:"syslog_drain_url"`
		AppURL             string      `json:"app_url"`
		ServiceInstanceURL string      `json:"service_instance_url"`
	} `json:"entity"`
}
