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
- Navigate to the [Choreo console](https://console.choreo.dev/).
- Click "Sign Up" if you don't already have an account.
- Select your authentication method (GitHub, Google, Microsoft, or email).
- Follow onboarding steps to create and configure your organization.
- Choose an organization name that reflects your company identity and is URL-friendly.
- Once logged in for the first time, select the "Platform Engineer (PE)" perspective. If you already have an organization set up, you can select the "PE" perspective from the dropdown menu at the top-right corner of the Choreo Console.

#### 2. Explore Platform Engineering Perspective Overview
- Centralized Data Plane Management: Clearly shows active/inactive data planes and their status.
- Built-in Environments: Easily view multiple environments (development, production).
- Comprehensive Insights: Quickly access operational insights for proactive troubleshooting.
- Integrated CD Pipelines: Default pipelines visible, simplifying deployment management.
- Egress and Workflow Controls: Configure security and automation workflows efficiently.

#### 3. Create a New Environment
- From the left-hand menu, click on Infrastructure, then select Environments.
- Examine existing environments (Development, Production) and their configurations.
- Click + Create Environment.
- Provide details for the new environment:
  - Name: Staging
  - Data Plane: e.g., Choreo Cloud US Dataplane
  - DNS Prefix: staging
- Review the automatically generated DNS URL.
- Decide if it's a production environment (leave unchecked for Staging).
- Click Create to provision the environment.
- Verify the new environment appears in your environment list.

#### 4. Update the Default CD Pipeline to Include the Staging Environment
- Navigate to DevOps → CD Pipelines from the Choreo navigation menu.
- Examine the existing CD Pipeline configuration, typically showing deployment from Development directly to Production.
- Click on the edit (✏️) icon next to the existing pipeline, such as "Default US Deployment Pipeline."
- Click on the "+" symbol between Development and Production environments.
- From the dropdown, select the newly created Staging environment.
- Review and confirm the pipeline sequence:
  - Development → Staging → Production
- Click Update to save your changes.
- Verify the pipeline now correctly shows deployment flows from Development → Staging → Production.

#### 5. Manage Team Access
- Navigate to Project Settings > Team.
- Invite "Joe" as a team member with "Project Admin" and "Developer" roles.
- Configure appropriate role-based permissions.

#### 6. Enable Environment Promotion Workflow
- Navigate to Governance → Workflows from the Choreo navigation menu.
- Locate the workflow named Environment Promotion.
- Toggle the Status switch to enable the workflow.
- Under Roles, select roles responsible for approving environment promotions (e.g., Admin, Project Admin, Choreo DevOps).
- Click Save to finalize your workflow configuration.
- Confirm that the workflow status is now enabled (green toggle).

#### 7. Provision Resources
- Navigate to DB & Services → Databases.
- Click + Create.
- Select PostgreSQL as the database type.
- Enter a meaningful service name (e.g., NonProductionPSQLServer).
- Click Next.
- Choose your preferred cloud provider and select a region close to your development team (e.g., United States).
- Select a suitable service plan for development purposes (e.g., Hobbyist).
- Review the configuration details and click Create to provision the database.
- Wait until the status changes from Creating to Active.
- Securely copy essential connection details (Host URL, Port, Default User, Default Database Name, Password).
- Securely store and share connection details with developers.

#### 8. Configure Security Policies
- Navigate to Governance → Egress Control from the Choreo navigation menu.
- Click + Create to set up a new egress policy.
- Select Deny All to block all egress traffic by default.
- Add an allow rule named Allow PSQL Server Access.
- Enter the PostgreSQL database server destination (host).
- Click Add Rule and confirm the rule appears correctly.
- Verify the policy is active and restricts outbound traffic only to the allowed PostgreSQL server.

#### 9. Invite Team Members
- Invite Team Members to the Developer Portal
  - Navigate to User Management → Users from the navigation panel.
  - Click + Invite Users at the top-right corner.
  - Enter email addresses for developers and relevant stakeholders.
  - Assign appropriate roles:
    - Developers: select "Developer" group.
    - Project Admins: select "Project Admin" group.
    - Architects: select or create the "Architect" group if necessary.
  - Verify emails and selected groups, then click Invite.
  - Confirm invitations under the Pending Invitations tab.
  - Invited team members will receive emails with instructions to join and start working.

</details>

<details>
<summary><h3>Part 2: Developer's Perspective</h3></summary>

#### 1. Access and Orientation
- Accept the invitation to join the Choreo project
- Complete onboarding steps:
  - Set up account and preferences
  - Install required CLI tools
  - Configure local development environment
- Explore the pre-configured environments and resources:
  - Review the development environment
  - Examine the staging environment
  - Understand the pipeline configuration
- Review documentation provided by the Platform Engineer

#### 2. Application Development
- Clone the application repository:
  - `git clone <repository-url>`
  - Set up local development environment
- Use the database connection details provided by the Platform Engineer:
  - Configure application to connect to the database
  - Set up appropriate connection pooling
  - Implement data access layer
- Implement key features:
  - User authentication and authorization
  - Customer data management
  - Reporting functionality

#### 3. Deployment and Testing
- Push code changes to trigger the CI/CD pipeline:
  - Commit and push changes to feature branch
  - Create pull request for review
  - Merge changes to main branch after approval
- Monitor the build and deployment process:
  - Track pipeline execution in Choreo dashboard
  - Review build logs and test results
  - Address any pipeline failures
- Verify application functionality in the staging environment:
  - Test key workflows
  - Validate database interactions
  - Ensure security policies are not violated

#### 4. Iterative Development
- Make code changes based on testing feedback:
  - Fix identified bugs
  - Refine user interfaces
  - Optimize performance
- Experience the streamlined development workflow:
  - Focus on code without infrastructure concerns
  - Rely on automated pipelines for testing and deployment
  - Leverage pre-configured environments
- Iterate rapidly on application features:
  - Implement new capabilities
  - Refine existing functionality
  - Improve user experience

</details>

<details>
<summary><h3>Part 3: Operational Excellence</h3></summary>

#### 1. Monitoring and Observability
- Platform Engineer navigates to the **Observability** section
- Set up comprehensive monitoring for:
  - Application health metrics:
    - Endpoint availability
    - Response times
    - Error rates
  - Resource utilization:
    - CPU and memory usage
    - Database connection pool status
    - Storage consumption
  - Business metrics:
    - Transaction volumes
    - User engagement
    - Feature usage

#### 2. Alert Configuration
- Create proactive alerts for:
  - Pipeline issues:
    - Failed deployments
    - Test failures
    - Build errors
  - Resource threshold breaches:
    - CPU utilization > 80%
    - Memory consumption > 85%
    - Database connections > 75%
  - Application-specific alerts:
    - Error rate spikes
    - Slow response times
    - Authentication failures
  - Configure notification channels:
    - Email alerts for critical issues
    - Slack notifications for team awareness
    - SMS alerts for urgent problems

#### 3. Continuous Improvement
- Review metrics and logs to identify optimization opportunities:
  - Analyze performance bottlenecks
  - Identify resource inefficiencies
  - Detect potential security issues
- Refine environment configurations:
  - Adjust resource allocations based on actual usage
  - Update environment variables for better performance
  - Optimize database settings
- Update pipeline steps for improved efficiency:
  - Parallelize compatible build steps
  - Optimize test execution
  - Implement caching strategies

</details>

<details>
<summary><h3>Part 4: Recap and Path Forward</h3></summary>

#### 1. Review of Achievements
- Demonstrate the completed application:
  - Showcase key features
  - Highlight performance metrics
  - Demonstrate scalability
- Analyze the benefits realized:
  - Reduced time-to-market for features
  - Improved code quality
  - Enhanced operational reliability
- Discuss the separation of concerns:
  - Platform Engineer focus on infrastructure
  - Developer focus on application code
  - Reduction in cross-domain context switching

</details>