# Building Applications with Choreo: An Internal Developer Platform Lab Session

## How an Internal Developer Platform Lets Developers Focus on Code

This hands-on lab session demonstrates how Choreo functions as an Internal Developer Platform (IDP), creating a clear separation of concerns between platform engineers and application developers. By implementing this separation, organizations can enable developers to focus entirely on writing code while platform engineers handle infrastructure, security, and governance concerns.

## Lab Outline

### Part 1: Platform Engineering Setup
Platform Engineers establish the foundational elements necessary for application development, including environments, pipelines, databases, and security policies.

### Part 2: Developer Experience
Developers leverage these pre-configured resources to build, test, and deploy applications with minimal infrastructure concerns.

### Part 3: Operational Excellence
Platform Engineers implement monitoring, alerting, and observability to ensure application reliability and performance.

### Part 4: Recap and Path Forward
Summary of achievements, benefits realized, and exploration of advanced capabilities for future implementation.

---

## Detailed Lab Instructions

<details>
<summary><h3>Part 1: Platform Engineer's Perspective</h3></summary>

#### 1. Sign Up to Choreo
- Navigate to the [Choreo console](https://console.choreo.dev/) and sign up using GitHub, Google, Microsoft, or email options.
- Complete the organization creation process by verifying the emails and OTP received.
- Select the "Platform Engineer (PE)" perspective from the Choreo Console.

#### 2. Explore Platform Engineering Perspective Overview
- Review the Data Plane Management section showing data plane status and regional distribution.
- Examine the pre-configured environments (Development, Production) with their specific runtime characteristics.
- Inspect the operational metrics dashboard showing resource utilization, deployment frequency, and system health.
- Verify the CD pipeline configurations that define the promotion paths between environments.
- Review network controls for managing ingress/egress traffic and security boundaries.

#### 3. Create a New Environment
- Navigate to Infrastructure → Environments and review the existing Development and Production environments.
- Click + Create Environment and enter the following configuration:
  - Name: "Staging"
  - Data Plane: "Choreo Cloud US Dataplane"
  - DNS Prefix: "staging"
  - Production Environment: Unchecked
- Review the auto-generated DNS URL: "staging.{org-name}.choreoapis.dev".
- Click Create and verify the Staging environment appears with "Active" status in your environment list.

#### 4. Update the Default CD Pipeline
- Navigate to DevOps → CD Pipelines and select the "Default US Deployment Pipeline".
- Click the edit (✏️) icon next to the pipeline.
- Click the "+" symbol between Development and Production environments.
- Select "Staging" from the environment dropdown.
- Review the pipeline sequence: Development → Staging → Production.
- Click Update to save your changes.
- Verify the pipeline visualization now shows the three-environment flow with arrows between each stage.

#### 5. Manage Team Access
- Navigate to Project Settings → Team.
- Click "Invite Member" and add user "Joe" with email "joe@example.com".
- Assign roles by checking "Project Admin" and "Developer" checkboxes.
- Set environment access permissions for Development, Staging, and Production.
- Click Send to issue the invitation.

#### 6. Enable Environment Promotion Workflow
- Navigate to Governance → Workflows.
- Find the "Environment Promotion" workflow in the list.
- Toggle the Status switch from Off to On.
- Under Roles, check "Admin", "Project Admin", and "Choreo DevOps".
- Click Save to apply the workflow configuration.
- Verify the workflow status shows Enabled with a green indicator.

#### 7. Provision Resources
- Navigate to DB & Services → Databases.
- Click + Create and select "PostgreSQL" as the database type.
- Configure the database:
  - Service name: "customer-portal-db"
  - Cloud provider: "Azure" or "AWS"
  - Region: "US East" or region closest to your location
  - Service plan: "Hobbyist" (1vCPU, 2GB RAM)
- Click Create and wait for status to change from "Creating" to "Active".
- Copy connection parameters:
  - Host: customer-portal-db.postgres.database.azure.com
  - Port: 5432
  - Default User: postgres
  - Default Database: postgres
  - Password: [auto-generated password]

#### 8. Configure Network Policies for Egress Control
- Navigate to Governance → Egress Control.
- Click + Create to establish a new policy.
- Select "Deny All" as the default rule.
- Add an allow rule with:
  - Name: "DB-Access"
  - Type: "Egress"
  - Target: "customer-portal-db.postgres.database.azure.com"
  - Port: 5432
  - Protocol: TCP
- Click Add Rule and verify the rule appears in the policy list.
- Confirm the policy status shows "Active" and test connectivity to confirm only database traffic is allowed.


</details>

<details>
<summary><h3>Part 2: Developer's Perspective</h3></summary>

#### 1. Access and Orientation
- Accept the project invitation email from Choreo and complete account setup.
- Install the Choreo CLI with: `npm install -g @choreo/cli` or `brew install choreo-cli`.
- Log into Choreo from CLI: `choreo login`.
- Explore the environments (Development, Staging, Production) from the Developer Console.
- Review the pipeline configuration and understand the promotion workflow.

#### 2. Application Development
- Clone the sample repository: `git clone https://github.com/choreo-samples/customer-portal.git`.
- Set up local environment with Node.js v14+ and npm v6+.
- Configure database connection in the application using provided credentials:
  ```javascript
  const db = {
    host: 'customer-portal-db.postgres.database.azure.com',
    port: 5432,
    database: 'postgres',
    user: 'postgres',
    password: process.env.DB_PASSWORD
  };
```

#### 3. Deployment and Testing
- Push code changes to trigger the CI/CD pipeline through the feature branch workflow and pull requests.
- Monitor build process, review logs and test results, and address any failures.
- Verify application functionality in the staging environment through comprehensive testing.

#### 4. Iterative Development
- Make code improvements based on testing feedback to fix bugs and optimize performance.
- Focus on code development without infrastructure concerns, leveraging the pre-configured pipelines.
- Rapidly implement and refine features to improve the overall user experience.

</details>

<details>
<summary><h3>Part 3: Operational Excellence</h3></summary>

#### 1. Monitoring and Observability
- Navigate to the Observability section and set up monitoring for application health, resource utilization, and business metrics.
- Configure comprehensive dashboards to track key performance indicators across all environments.

#### 2. Alert Configuration
- Create alerts for pipeline issues, resource threshold breaches, and application-specific problems.
- Configure appropriate notification channels (email, Slack, SMS) based on issue severity.

#### 3. Continuous Improvement
- Analyze metrics and logs to identify bottlenecks, inefficiencies, and potential issues.
- Refine environment configurations and optimize pipeline steps for better performance.

</details>

<details>
<summary><h3>Part 4: Recap and Path Forward</h3></summary>

#### 1. Review of Achievements
- Demonstrate the completed application highlighting key features and performance metrics.
- Analyze benefits including faster development cycles, improved code quality, and operational reliability.
- Discuss how separation of concerns between platform engineering and development improved productivity.

#### 2. Advanced Capabilities
- Explore advanced deployment strategies, security configurations, integration capabilities, and enhanced observability tools for future implementation.

</details>