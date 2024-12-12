package ledger

import (
	"fmt"
	"time"
)

// DeleteAPIProxy removes an API proxy configuration
func (l *IntegrationLedger) DeleteAPIProxy(proxyID string) error {
	if _, exists := l.APIProxies[proxyID]; !exists {
		return fmt.Errorf("API proxy with ID %s does not exist", proxyID)
	}
	delete(l.APIProxies, proxyID)
	return nil
}

// ApplyUpdate applies an encrypted application update
func (l *IntegrationLedger) ApplyUpdate(appID string, encryptedUpdate ApplicationUpdate) error {
	if _, exists := l.Applications[appID]; !exists {
		return fmt.Errorf("application with ID %s does not exist", appID)
	}
	l.Applications[appID] = encryptedUpdate
	return nil
}

// RevertUpdate reverts an application to a previous version
func (l *IntegrationLedger) RevertUpdate(appID string, version string) error {
	app, exists := l.Applications[appID]
	if !exists {
		return fmt.Errorf("application with ID %s does not exist", appID)
	}
	if app.Version != version {
		return fmt.Errorf("no rollback available for version %s", version)
	}
	// Perform rollback logic (e.g., restore backup)
	return nil
}

// AddServiceProvider adds a new service provider
func (l *IntegrationLedger) AddServiceProvider(provider ServiceProvider) error {
	if _, exists := l.ServiceProviders[provider.ProviderID]; exists {
		return fmt.Errorf("service provider with ID %s already exists", provider.ProviderID)
	}
	l.ServiceProviders[provider.ProviderID] = provider
	return nil
}

// RemoveServiceProvider removes an existing service provider
func (l *IntegrationLedger) RemoveServiceProvider(providerID string) error {
	if _, exists := l.ServiceProviders[providerID]; !exists {
		return fmt.Errorf("service provider with ID %s does not exist", providerID)
	}
	delete(l.ServiceProviders, providerID)
	return nil
}

// GetIntegrationState retrieves the state of a proxy integration
func (l *IntegrationLedger) GetIntegrationState(proxyID string) (IntegrationState, error) {
	state, exists := l.IntegrationStates[proxyID]
	if !exists {
		return IntegrationState{}, fmt.Errorf("integration state for proxy ID %s not found", proxyID)
	}
	return state, nil
}

// SetVersionControl enables or disables version control for an application
func (l *IntegrationLedger) SetVersionControl(appID string, enabled bool) error {
	l.VersionControlEnabled[appID] = enabled
	return nil
}

// UpdateServiceConfig updates a service's configuration.
func (l *IntegrationLedger) UpdateServiceConfig(serviceID string, config ServiceConfig) error {
	l.ServiceConfigs[serviceID] = config
	return nil
}

// AddWebhook adds a webhook configuration for a service.
func (l *IntegrationLedger) AddWebhook(serviceID string, webhook WebhookConfig) error {
	l.Webhooks[serviceID] = append(l.Webhooks[serviceID], webhook)
	return nil
}

// DeleteWebhook removes a webhook configuration for a service.
func (l *IntegrationLedger) DeleteWebhook(serviceID string, webhookID string) error {
	webhooks, exists := l.Webhooks[serviceID]
	if !exists {
		return fmt.Errorf("service ID %s has no webhooks", serviceID)
	}
	var updatedWebhooks []WebhookConfig
	for _, webhook := range webhooks {
		if webhook.WebhookID != webhookID {
			updatedWebhooks = append(updatedWebhooks, webhook)
		}
	}
	if len(updatedWebhooks) == len(webhooks) {
		return fmt.Errorf("webhook ID %s not found for service ID %s", webhookID, serviceID)
	}
	l.Webhooks[serviceID] = updatedWebhooks
	return nil
}

// AddFunctionToService adds a custom function to a service.
func (l *IntegrationLedger) AddFunctionToService(serviceID string, function CustomFunction) error {
	l.CustomFunctions[serviceID] = append(l.CustomFunctions[serviceID], function)
	return nil
}

// RemoveFunctionFromService removes a custom function from a service.
func (l *IntegrationLedger) RemoveFunctionFromService(serviceID string, functionID string) error {
	functions, exists := l.CustomFunctions[serviceID]
	if !exists {
		return fmt.Errorf("service ID %s has no custom functions", serviceID)
	}
	var updatedFunctions []CustomFunction
	for _, function := range functions {
		if function.FunctionID != functionID {
			updatedFunctions = append(updatedFunctions, function)
		}
	}
	if len(updatedFunctions) == len(functions) {
		return fmt.Errorf("function ID %s not found for service ID %s", functionID, serviceID)
	}
	l.CustomFunctions[serviceID] = updatedFunctions
	return nil
}

// AddAnalyticsToService integrates an analytics tool with a service.
func (l *IntegrationLedger) AddAnalyticsToService(serviceID string, analytics AnalyticsConfig) error {
	l.AnalyticsTools[serviceID] = append(l.AnalyticsTools[serviceID], analytics)
	return nil
}

// RemoveAnalyticsFromService removes an analytics tool from a service.
func (l *IntegrationLedger) RemoveAnalyticsFromService(serviceID string, analyticsID string) error {
	analytics, exists := l.AnalyticsTools[serviceID]
	if !exists {
		return fmt.Errorf("service ID %s has no analytics tools", serviceID)
	}
	var updatedAnalytics []AnalyticsConfig
	for _, tool := range analytics {
		if tool.AnalyticsID != analyticsID {
			updatedAnalytics = append(updatedAnalytics, tool)
		}
	}
	if len(updatedAnalytics) == len(analytics) {
		return fmt.Errorf("analytics ID %s not found for service ID %s", analyticsID, serviceID)
	}
	l.AnalyticsTools[serviceID] = updatedAnalytics
	return nil
}

// UpdateAppConfig updates an application's configuration.
func (l *IntegrationLedger) UpdateAppConfig(appID string, config AppConfig) error {
	l.AppConfigs[appID] = config
	return nil
}

// SetAppStatus sets the active/inactive status of an application.
func (l *IntegrationLedger) SetAppStatus(appID string, status bool) error {
	l.AppStatus[appID] = status
	return nil
}

// VerifyAPIVersionCompatibility checks API compatibility for a specific application.
func (l *IntegrationLedger) VerifyAPIVersionCompatibility(appID, apiVersion string) (bool, error) {
	if versions, exists := l.APICompatibility[appID]; exists {
		if isCompatible, ok := versions[apiVersion]; ok {
			return isCompatible, nil
		}
	}
	return false, fmt.Errorf("API version %s not compatible with app %s", apiVersion, appID)
}

// SetIntegrationParameters sets the integration parameters for an application.
func (l *IntegrationLedger) SetIntegrationParameters(appID string, params IntegrationParams) error {
	l.IntegrationParameters[appID] = params
	return nil
}

// ValidateSchema validates the API schema for a given application.
func (l *IntegrationLedger) ValidateSchema(appID string, schema APISchema) (bool, error) {
	storedSchema, exists := l.APISchemas[appID]
	if !exists {
		return false, fmt.Errorf("no schema found for app %s", appID)
	}
	return storedSchema.Definition == schema.Definition, nil
}

// AddExtension adds an extension to an application.
func (l *IntegrationLedger) AddExtension(appID string, extension Extension) error {
	l.Extensions[appID] = append(l.Extensions[appID], extension)
	return nil
}

// RemoveExtension removes an extension from an application.
func (l *IntegrationLedger) RemoveExtension(appID, extensionID string) error {
	extensions, exists := l.Extensions[appID]
	if !exists {
		return fmt.Errorf("no extensions found for app %s", appID)
	}
	var updatedExtensions []Extension
	for _, ext := range extensions {
		if ext.ExtensionID != extensionID {
			updatedExtensions = append(updatedExtensions, ext)
		}
	}
	if len(updatedExtensions) == len(extensions) {
		return fmt.Errorf("extension ID %s not found for app %s", extensionID, appID)
	}
	l.Extensions[appID] = updatedExtensions
	return nil
}

// AddEventHandler registers an event handler for an application.
func (l *IntegrationLedger) AddEventHandler(appID string, handler EventHandler) error {
	l.EventHandlers[appID] = append(l.EventHandlers[appID], handler)
	return nil
}

// RemoveEventHandler removes an event handler for a specific application.
func (l *IntegrationLedger) RemoveEventHandler(appID, handlerID string) error {
	handlers, exists := l.EventHandlers[appID]
	if !exists {
		return fmt.Errorf("no event handlers found for app %s", appID)
	}
	var updatedHandlers []EventHandler
	for _, handler := range handlers {
		if handler.HandlerID != handlerID {
			updatedHandlers = append(updatedHandlers, handler)
		}
	}
	if len(updatedHandlers) == len(handlers) {
		return fmt.Errorf("event handler %s not found for app %s", handlerID, appID)
	}
	l.EventHandlers[appID] = updatedHandlers
	return nil
}

// AddLibrary integrates a library into an application.
func (l *IntegrationLedger) AddLibrary(appID string, library Library) error {
	l.Libraries[appID] = append(l.Libraries[appID], library)
	return nil
}

// RemoveLibrary removes a library from an application.
func (l *IntegrationLedger) RemoveLibrary(appID, libraryID string) error {
	libraries, exists := l.Libraries[appID]
	if !exists {
		return fmt.Errorf("no libraries found for app %s", appID)
	}
	var updatedLibraries []Library
	for _, lib := range libraries {
		if lib.LibraryID != libraryID {
			updatedLibraries = append(updatedLibraries, lib)
		}
	}
	if len(updatedLibraries) == len(libraries) {
		return fmt.Errorf("library %s not found for app %s", libraryID, appID)
	}
	l.Libraries[appID] = updatedLibraries
	return nil
}

// UpdateAPIKeys updates the API keys for a specific application.
func (l *IntegrationLedger) UpdateAPIKeys(appID string, apiKeys APIKeys) error {
	l.APIKeys[appID] = apiKeys
	return nil
}

// CheckIntegrationStatus checks the integration status of an application.
func (l *IntegrationLedger) CheckIntegrationStatus(appID string) (IntegrationStatus, error) {
	status, exists := l.IntegrationStatuses[appID]
	if !exists {
		return IntegrationStatus{}, fmt.Errorf("integration status not found for app %s", appID)
	}
	return status, nil
}

// ActivateServiceIntegration enables a service for a specific application.
func (l *IntegrationLedger) ActivateServiceIntegration(appID string, service Service) error {
	l.ServiceIntegrations[appID] = append(l.ServiceIntegrations[appID], service)
	return nil
}

// AddDapp adds a new DApp to the ledger.
func (l *IntegrationLedger) AddDapp(metadata DappMetadata) error {
	if _, exists := l.Dapps[metadata.DappID]; exists {
		return fmt.Errorf("DApp with ID %s already exists", metadata.DappID)
	}
	l.Dapps[metadata.DappID] = metadata
	return nil
}

// UpdateDapp updates an existing DApp in the ledger.
func (l *IntegrationLedger) UpdateDapp(dappID string, metadata DappMetadata) error {
	if _, exists := l.Dapps[dappID]; !exists {
		return fmt.Errorf("DApp with ID %s does not exist", dappID)
	}
	l.Dapps[dappID] = metadata
	return nil
}

// DeleteDapp removes a DApp from the ledger.
func (l *IntegrationLedger) DeleteDapp(dappID string) error {
	if _, exists := l.Dapps[dappID]; !exists {
		return fmt.Errorf("DApp with ID %s does not exist", dappID)
	}
	delete(l.Dapps, dappID)
	return nil
}

// AddAPIEndpoint adds an API endpoint to a DApp.
func (l *IntegrationLedger) AddAPIEndpoint(dappID string, endpoint APIEndpoint) error {
	if _, exists := l.Dapps[dappID]; !exists {
		return fmt.Errorf("DApp with ID %s does not exist", dappID)
	}
	l.APIEndpoints[dappID] = append(l.APIEndpoints[dappID], endpoint)
	return nil
}

// DeleteAPIEndpoint removes an API endpoint from a DApp.
func (l *IntegrationLedger) DeleteAPIEndpoint(dappID, endpointID string) error {
	endpoints, exists := l.APIEndpoints[dappID]
	if !exists {
		return fmt.Errorf("no API endpoints found for DApp %s", dappID)
	}
	var updatedEndpoints []APIEndpoint
	for _, ep := range endpoints {
		if ep.EndpointID != endpointID {
			updatedEndpoints = append(updatedEndpoints, ep)
		}
	}
	if len(updatedEndpoints) == len(endpoints) {
		return fmt.Errorf("API endpoint %s not found for DApp %s", endpointID, dappID)
	}
	l.APIEndpoints[dappID] = updatedEndpoints
	return nil
}

// AddCLICommand adds a CLI command to a DApp.
func (l *IntegrationLedger) AddCLICommand(dappID string, command CLICommand) error {
	if _, exists := l.Dapps[dappID]; !exists {
		return fmt.Errorf("DApp with ID %s does not exist", dappID)
	}
	l.CLICommands[dappID] = append(l.CLICommands[dappID], command)
	return nil
}

// DeleteCLICommand removes a CLI command from a DApp.
func (l *IntegrationLedger) DeleteCLICommand(dappID, commandID string) error {
	commands, exists := l.CLICommands[dappID]
	if !exists {
		return fmt.Errorf("no CLI commands found for DApp %s", dappID)
	}
	var updatedCommands []CLICommand
	for _, cmd := range commands {
		if cmd.CommandID != commandID {
			updatedCommands = append(updatedCommands, cmd)
		}
	}
	if len(updatedCommands) == len(commands) {
		return fmt.Errorf("CLI command %s not found for DApp %s", commandID, dappID)
	}
	l.CLICommands[dappID] = updatedCommands
	return nil
}

// VerifyFunctionality validates new functionality added to a DApp.
func (l *IntegrationLedger) VerifyFunctionality(dappID string, functionality Functionality) (bool, error) {
	if _, exists := l.Dapps[dappID]; !exists {
		return false, fmt.Errorf("DApp with ID %s does not exist", dappID)
	}
	// Placeholder for actual security checks
	if functionality.Name == "" || functionality.Code == "" {
		return false, fmt.Errorf("invalid functionality data")
	}
	return true, nil
}

// ToggleFeature toggles a feature for a DApp.
func (l *IntegrationLedger) ToggleFeature(dappID, featureName string, enabled bool) error {
	if _, exists := l.FeatureToggles[dappID]; !exists {
		l.FeatureToggles[dappID] = []FeatureToggle{}
	}
	for i, toggle := range l.FeatureToggles[dappID] {
		if toggle.FeatureName == featureName {
			l.FeatureToggles[dappID][i].Enabled = enabled
			l.FeatureToggles[dappID][i].Timestamp = time.Now()
			return nil
		}
	}
	l.FeatureToggles[dappID] = append(l.FeatureToggles[dappID], FeatureToggle{
		DappID:      dappID,
		FeatureName: featureName,
		Enabled:     enabled,
		Timestamp:   time.Now(),
	})
	return nil
}

// GetIntegrationStatus retrieves the integration status of a DApp or application.
func (l *IntegrationLedger) GetIntegrationStatus(dappID string) (IntegrationStatus, error) {
	l.Lock()
	defer l.Unlock()

	// Check if the integration status exists for the given dappID
	status, exists := l.IntegrationStatuses[dappID]
	if !exists {
		return IntegrationStatus{}, fmt.Errorf("integration status not found for DApp ID: %s", dappID)
	}

	return status, nil
}


// AddExternalService adds an external service to a DApp.
func (l *IntegrationLedger) AddExternalService(dappID string, service ExternalService) error {
	if _, exists := l.ExternalServices[dappID]; !exists {
		l.ExternalServices[dappID] = []ExternalService{}
	}
	l.ExternalServices[dappID] = append(l.ExternalServices[dappID], service)
	return nil
}

// DeleteExternalService removes an external service from a DApp.
func (l *IntegrationLedger) DeleteExternalService(dappID, serviceID string) error {
	services, exists := l.ExternalServices[dappID]
	if !exists {
		return fmt.Errorf("no services found for DApp %s", dappID)
	}
	var updatedServices []ExternalService
	for _, service := range services {
		if service.ServiceID != serviceID {
			updatedServices = append(updatedServices, service)
		}
	}
	if len(updatedServices) == len(services) {
		return fmt.Errorf("service %s not found for DApp %s", serviceID, dappID)
	}
	l.ExternalServices[dappID] = updatedServices
	return nil
}

// AddOpcode adds a new opcode to a DApp.
func (l *IntegrationLedger) AddOpcode(dappID string, opcode Opcode) error {
	if _, exists := l.Opcodes[dappID]; !exists {
		l.Opcodes[dappID] = []Opcode{}
	}
	l.Opcodes[dappID] = append(l.Opcodes[dappID], opcode)
	return nil
}

// RemoveOpcode removes an opcode from a DApp.
func (l *IntegrationLedger) RemoveOpcode(dappID, opcodeID string) error {
	opcodes, exists := l.Opcodes[dappID]
	if !exists {
		return fmt.Errorf("no opcodes found for DApp %s", dappID)
	}
	var updatedOpcodes []Opcode
	for _, opcode := range opcodes {
		if opcode.OpcodeID != opcodeID {
			updatedOpcodes = append(updatedOpcodes, opcode)
		}
	}
	if len(updatedOpcodes) == len(opcodes) {
		return fmt.Errorf("opcode %s not found for DApp %s", opcodeID, dappID)
	}
	l.Opcodes[dappID] = updatedOpcodes
	return nil
}

// AddAppComponent adds an application component to a DApp.
func (l *IntegrationLedger) AddAppComponent(dappID string, component AppComponent) error {
	if _, exists := l.AppComponents[dappID]; !exists {
		l.AppComponents[dappID] = []AppComponent{}
	}
	l.AppComponents[dappID] = append(l.AppComponents[dappID], component)
	return nil
}

// CheckFeatureDependencies validates that all dependencies for a feature are met.
func (l *IntegrationLedger) CheckFeatureDependencies(featureID string, dependencies []string) (bool, error) {
	for _, dep := range dependencies {
		found := false
		for _, existingDep := range l.FeatureDependencies[featureID] {
			if existingDep.DependencyID == dep {
				found = true
				break
			}
		}
		if !found {
			return false, fmt.Errorf("missing dependency: %s", dep)
		}
	}
	return true, nil
}

// GetApplicationFeatures retrieves all features for a given application.
func (l *IntegrationLedger) GetApplicationFeatures(appID string) ([]Feature, error) {
	features, exists := l.ApplicationFeatures[appID]
	if !exists {
		return nil, fmt.Errorf("no features found for application %s", appID)
	}
	return features, nil
}

// AddFeatureDependency adds a dependency to a feature.
func (l *IntegrationLedger) AddFeatureDependency(featureID string, dependency Dependency) error {
	if _, exists := l.FeatureDependencies[featureID]; !exists {
		l.FeatureDependencies[featureID] = []Dependency{}
	}
	l.FeatureDependencies[featureID] = append(l.FeatureDependencies[featureID], dependency)
	return nil
}

// RemoveFeatureDependency removes a specific dependency from a feature.
func (l *IntegrationLedger) RemoveFeatureDependency(featureID, dependencyID string) error {
	dependencies, exists := l.FeatureDependencies[featureID]
	if !exists {
		return fmt.Errorf("no dependencies found for feature %s", featureID)
	}
	var updatedDeps []Dependency
	for _, dep := range dependencies {
		if dep.DependencyID != dependencyID {
			updatedDeps = append(updatedDeps, dep)
		}
	}
	if len(updatedDeps) == len(dependencies) {
		return fmt.Errorf("dependency %s not found for feature %s", dependencyID, featureID)
	}
	l.FeatureDependencies[featureID] = updatedDeps
	return nil
}

// SetAPIGatewayStatus toggles the API Gateway status for an application.
func (l *IntegrationLedger) SetAPIGatewayStatus(appID string, status bool) error {
	l.APIGateways[appID] = status
	return nil
}

// UpdateIntegrationMapping updates cross-application integration mappings.
func (l *IntegrationLedger) UpdateIntegrationMapping(appID string, mapping IntegrationMapping) error {
	l.IntegrationMappings[appID] = mapping
	return nil
}

// ConfigureWorkflow configures an integration workflow for a given application.
func (l *IntegrationLedger) ConfigureWorkflow(appID string, workflow WorkflowConfig) error {
	l.Workflows[appID] = workflow
	return nil
}

// ReviewSecurity reviews the security of an application's integration.
func (l *IntegrationLedger) ReviewSecurity(appID string) (SecurityReview, error) {
	review, exists := l.SecurityReviews[appID]
	if !exists {
		return SecurityReview{}, fmt.Errorf("no security review found for application %s", appID)
	}
	return review, nil
}

// LogIntegrationActivity logs an integration activity for traceability.
func (l *IntegrationLedger) LogIntegrationActivity(appID string, activity ActivityLog) error {
	if _, exists := l.IntegrationActivities[appID]; !exists {
		l.IntegrationActivities[appID] = []ActivityLog{}
	}
	l.IntegrationActivities[appID] = append(l.IntegrationActivities[appID], activity)
	return nil
}

// CheckComponentCompatibility verifies compatibility between two components.
func (l *IntegrationLedger) CheckComponentCompatibility(appID, componentID string) (bool, error) {
	components, exists := l.ComponentCompatibility[appID]
	if !exists {
		return false, fmt.Errorf("no compatibility data found for application %s", appID)
	}
	compatible, found := components[componentID]
	if !found {
		return false, fmt.Errorf("component %s not found in compatibility data for application %s", componentID, appID)
	}
	return compatible, nil
}

// ExecuteCrossAppFunction executes a cross-application function.
func (l *IntegrationLedger) ExecuteCrossAppFunction(sourceAppID, targetAppID string, function CrossAppFunction) error {
	functionID := fmt.Sprintf("%s:%s:%s", sourceAppID, targetAppID, function.FunctionID)
	l.CrossAppFunctions[functionID] = function
	return nil
}

// AddDependentModule registers a module dependency for a feature.
func (l *IntegrationLedger) AddDependentModule(featureID string, module Module) error {
	if _, exists := l.DependentModules[featureID]; !exists {
		l.DependentModules[featureID] = []Module{}
	}
	l.DependentModules[featureID] = append(l.DependentModules[featureID], module)
	return nil
}

// RemoveDependentModule removes a module dependency from a feature.
func (l *IntegrationLedger) RemoveDependentModule(featureID, moduleID string) error {
	modules, exists := l.DependentModules[featureID]
	if !exists {
		return fmt.Errorf("no modules found for feature %s", featureID)
	}
	var updatedModules []Module
	for _, module := range modules {
		if module.ModuleID != moduleID {
			updatedModules = append(updatedModules, module)
		}
	}
	if len(updatedModules) == len(modules) {
		return fmt.Errorf("module %s not found for feature %s", moduleID, featureID)
	}
	l.DependentModules[featureID] = updatedModules
	return nil
}

// SetServiceIntegrationStatus enables or disables a service integration.
func (l *IntegrationLedger) SetServiceIntegrationStatus(serviceID string, status bool) error {
	l.ServiceIntegrationStatus[serviceID] = status
	return nil
}

// StoreAPIResponse stores the encrypted response from an external API.
func (l *IntegrationLedger) StoreAPIResponse(serviceID string, response string) error {
	l.APIResponses[serviceID] = response
	return nil
}

// CheckServiceIntegration checks the integration status of a service.
func (l *IntegrationLedger) CheckServiceIntegration(serviceID string) (bool, error) {
	status, exists := l.ServiceIntegrationStatus[serviceID]
	if !exists {
		return false, fmt.Errorf("service integration status not found for %s", serviceID)
	}
	return status, nil
}

// CheckMethodCompatibility verifies if a method is compatible with a service.
func (l *Ledger) CheckMethodCompatibility(serviceID, methodName string) (bool, error) {
	// Placeholder for method compatibility logic.
	// Add logic for compatibility checks as needed.
	return true, nil
}

// AddPolicy adds a policy to the service.
func (l *IntegrationLedger) AddPolicy(serviceID string, policy Policy) error {
	if _, exists := l.ServicePolicies[serviceID]; !exists {
		l.ServicePolicies[serviceID] = []Policy{}
	}
	l.ServicePolicies[serviceID] = append(l.ServicePolicies[serviceID], policy)
	return nil
}

// RemovePolicy removes a policy from a service.
func (l *IntegrationLedger) RemovePolicy(serviceID, policyID string) error {
	policies, exists := l.ServicePolicies[serviceID]
	if !exists {
		return fmt.Errorf("no policies found for service %s", serviceID)
	}
	var updatedPolicies []Policy
	for _, policy := range policies {
		if policy.PolicyID != policyID {
			updatedPolicies = append(updatedPolicies, policy)
		}
	}
	if len(updatedPolicies) == len(policies) {
		return fmt.Errorf("policy %s not found for service %s", policyID, serviceID)
	}
	l.ServicePolicies[serviceID] = updatedPolicies
	return nil
}

// LogServiceEvent records an integration event.
func (l *IntegrationLedger) LogServiceEvent(serviceID string, event IntegrationEvent) error {
	if _, exists := l.IntegrationEvents[serviceID]; !exists {
		l.IntegrationEvents[serviceID] = []IntegrationEvent{}
	}
	l.IntegrationEvents[serviceID] = append(l.IntegrationEvents[serviceID], event)
	return nil
}

// SetAccessLevel defines access permissions for a service.
func (l *IntegrationLedger) SetAccessLevel(serviceID string, accessLevel AccessLevel) error {
	l.AccessLevels[serviceID] = accessLevel
	return nil
}

// UpdateLogLevel updates the logging level for a service.
func (l *IntegrationLedger) UpdateLogLevel(serviceID string, logLevel LogLevel) error {
	l.LogLevels[serviceID] = logLevel
	return nil
}

// GetIntegrationLogs fetches logs for a specific service integration.
func (l *IntegrationLedger) GetIntegrationLogs(serviceID string) ([]IntegrationLog, error) {
	logs, exists := l.IntegrationLogs[serviceID]
	if !exists {
		return nil, fmt.Errorf("no logs found for service %s", serviceID)
	}
	return logs, nil
}

// GetIntegrationHealthStatus retrieves the health status of a service integration.
func (l *IntegrationLedger) GetIntegrationHealthStatus(serviceID string) (HealthStatus, error) {
	status, exists := l.IntegrationHealth[serviceID]
	if !exists {
		return HealthStatus{}, fmt.Errorf("no health status found for service %s", serviceID)
	}
	return status, nil
}

// AttachExtension adds an extension to a DApp.
func (l *IntegrationLedger) AttachExtension(dappID string, extension Extension) error {
	if _, exists := l.DappExtensions[dappID]; !exists {
		l.DappExtensions[dappID] = []Extension{}
	}
	l.DappExtensions[dappID] = append(l.DappExtensions[dappID], extension)
	return nil
}

// DetachExtension removes an extension from a DApp.
func (l *IntegrationLedger) DetachExtension(dappID, extensionID string) error {
	extensions, exists := l.DappExtensions[dappID]
	if !exists {
		return fmt.Errorf("no extensions found for DApp %s", dappID)
	}
	var updatedExtensions []Extension
	for _, ext := range extensions {
		if ext.ExtensionID != extensionID {
			updatedExtensions = append(updatedExtensions, ext)
		}
	}
	l.DappExtensions[dappID] = updatedExtensions
	return nil
}

// RunIntegrationTest executes a test for a service integration.
func (l *IntegrationLedger) RunIntegrationTest(serviceID string, testConfig TestConfig) (bool, error) {
	tests, exists := l.IntegrationTests[serviceID]
	if !exists {
		l.IntegrationTests[serviceID] = []TestConfig{}
	}
	l.IntegrationTests[serviceID] = append(tests, testConfig)
	// Placeholder for test execution logic
	return true, nil
}

// ValidateCLITool verifies the compatibility of a CLI tool.
func (l *IntegrationLedger) ValidateCLITool(cliTool CLITool) (bool, error) {
	// Placeholder for validation logic
	if cliTool.Compatibility == "Compatible" {
		return true, nil
	}
	return false, nil
}

// InstallAPIProxy installs an API proxy for a service.
func (l *IntegrationLedger) InstallAPIProxy(serviceID string, proxyConfig APIProxyConfig) error {
	l.APIProxies[serviceID] = proxyConfig
	return nil
}
