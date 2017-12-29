// +build !dev

package graphqlbackend

// Code generated by schema_generate.go

// Schema is the raw graqhql schema
var Schema = `# Run this before committing changes to this file
# go generate sourcegraph.com/sourcegraph/sourcegraph/cmd/frontend/internal/graphqlbackend

schema {
	query: Query
	mutation: Mutation
}

type EmptyResponse {
	alwaysNil: String
}

interface Node {
	id: ID!
}

type ThreadLines {
	# HTML context lines before 'html'.
	#
	# It is sanitized already by the server, and thus is safe for rendering.
	htmlBefore(isLightTheme: Boolean!): String!
	# HTML lines that the user's selection was made on.
	#
	# It is sanitized already by the server, and thus is safe for rendering.
	html(isLightTheme: Boolean!): String!
	# HTML context lines after 'html'.
	#
	# It is sanitized already by the server, and thus is safe for rendering.
	htmlAfter(isLightTheme: Boolean!): String!
	# text context lines before 'text'.
	textBefore: String!
	# text lines that the user's selection was made on.
	text: String!
	# text context lines after 'text'.
	textAfter: String!
	# byte offset into textLines where user selection began
	#
	# In Go syntax, userSelection := text[rangeStart:rangeStart+rangeLength]
	textSelectionRangeStart: Int!
	# length in bytes of the user selection
	#
	# In Go syntax, userSelection := text[rangeStart:rangeStart+rangeLength]
	textSelectionRangeLength: Int!
}

# Literally the exact same thing as above, except it's an input type because
# GraphQL doesn't allow mixing input and output types.
input ThreadLinesInput {
	# HTML context lines before 'html'.
	htmlBefore: String!
	# HTML lines that the user's selection was made on.
	html: String!
	# HTML context lines after 'html'.
	htmlAfter: String!
	# text context lines before 'text'.
	textBefore: String!
	# text lines that the user's selection was made on.
	text: String!
	# text context lines after 'text'.
	textAfter: String!
	# byte offset into textLines where user selection began
	#
	# In Go syntax, userSelection := text[rangeStart:rangeStart+rangeLength]
	textSelectionRangeStart: Int!
	# length in bytes of the user selection
	#
	# In Go syntax, userSelection := text[rangeStart:rangeStart+rangeLength]
	textSelectionRangeLength: Int!
}

input CreateThreadInput {
	orgID: ID!
	canonicalRemoteID: String!
	cloneURL: String!
	repoRevisionPath: String!
	linesRevisionPath: String!
	repoRevision: String!
	linesRevision: String!
	branch: String
	startLine: Int!
	endLine: Int!
	startCharacter: Int!
	endCharacter: Int!
	rangeLength: Int!
	contents: String!
	lines: ThreadLinesInput
}

# A string containing valid JSON.
scalar JSONString

type Mutation {
	createThread(
		orgID: ID!
		canonicalRemoteID: String!
		cloneURL: String!
		file: String!
		repoRevision: String!
		linesRevision: String!
		branch: String
		startLine: Int!
		endLine: Int!
		startCharacter: Int!
		endCharacter: Int!
		rangeLength: Int!
		contents: String!
		lines: ThreadLinesInput
	): Thread! @deprecated(reason: "use createThread2")
	createThread2(input: CreateThreadInput!): Thread!
	updateUser(username: String, displayName: String, avatarURL: String): User!
	# Update the settings for the currently authenticated user.
	updateUserSettings(lastKnownSettingsID: Int, contents: String!): Settings!
	updateThread(threadID: Int!, archived: Boolean): Thread!
	addCommentToThread(threadID: Int!, contents: String!): Thread!
	# This method is the same as addCommentToThread, the only difference is
	# that authentication is based on the secret ULID instead of the current
	# user.
	#
	# 🚨 SECURITY: Every field of the return type here is accessible publicly
	# given a shared item URL.
	addCommentToThreadShared(ulid: String!, threadID: Int!, contents: String!): SharedItemThread!
	shareThread(threadID: Int!): String!
	shareComment(commentID: Int!): String!
	createOrg(name: String!, displayName: String!): Org!
	updateOrg(id: ID!, displayName: String, slackWebhookURL: String): Org!
	updateOrgSettings(
		# The ID of the org whose settings should be updated.
		id: ID
		# TODO(sqs): orgID is deprecated. Use id instead.
		orgID: ID
		lastKnownSettingsID: Int
		contents: String!
	): Settings!
	# Deletes an organization. Only site admins may perform this mutation.
	deleteOrganization(organization: ID!): EmptyResponse
	# Creates a user account for a new user and generates a reset password link that the user
	# must visit to sign into the account. Only site admins may perform this mutation.
	createUserBySiteAdmin(username: String!, email: String!): CreateUserBySiteAdminResult!
	# Randomize a user's password so that they need to reset it before they can sign in again.
	# Only site admins may perform this mutation.
	randomizeUserPasswordBySiteAdmin(user: ID!): RandomizeUserPasswordBySiteAdminResult!
	# Deletes a user account. Only site admins may perform this mutation.
	deleteUser(user: ID!): EmptyResponse
	inviteUser(email: String!, orgID: ID!): InviteUserResult
	acceptUserInvite(inviteToken: String!): OrgInviteStatus!
	removeUserFromOrg(userID: ID!, orgID: ID!): EmptyResponse
	# adds a phabricator repository to the Sourcegraph server.
	# example callsign: "MUX"
	# example uri: "github.com/gorilla/mux"
	addPhabricatorRepo(callsign: String!, uri: String!, url: String!): EmptyResponse
	logUserEvent(event: UserEvent!): EmptyResponse
	# All mutations that update configuration settings are under this field.
	configurationMutation(input: ConfigurationMutationGroupInput!): ConfigurationMutation
	# Updates the site configuration.
	updateSiteConfiguration(input: String!): EmptyResponse
	# Sets whether the user with the specified user ID is a site admin.
	#
	# 🚨 SECURITY: Only trusted users should be given site admin permissions.
	# Site admins have full access to the server's site configuration and other
	# sensitive data, and they can perform destructive actions such as
	# restarting the site.
	setUserIsSiteAdmin(userID: ID!, siteAdmin: Boolean!): EmptyResponse
	# Reloads the site by restarting the server. This is not supported for all deployment
	# types. This may cause downtime.
	reloadSite: EmptyResponse
}

# Input for Mutation.configuration, which contains fields that all configuration
# mutations need.
input ConfigurationMutationGroupInput {
	# The subject whose configuration to mutate (org, user, etc.).
	subject: ID!
	# The ID of the last-known configuration known to the client, or null if
	# there is none. This field is used to prevent race conditions when there
	# are concurrent editors.
	lastID: Int
}

# Mutations that update configuration settings. These mutations are grouped
# together because they:
#
# - are all versioned to avoid race conditions with concurrent editors
# - all apply to a specific configuration subject
#
# Grouping them lets us extract those common parameters to the
# Mutation.configuration field.
type ConfigurationMutation {
	# Perform a raw configuration update. Use one of the other fields on this
	# type instead if possible.
	updateConfiguration(input: UpdateConfigurationInput!): UpdateConfigurationPayload
	# Create a saved query.
	createSavedQuery(description: String!, query: String!, scopeQuery: String!): SavedQuery!
	# Update the saved query with the given ID in the configuration.
	updateSavedQuery(id: ID!, description: String, query: String, scopeQuery: String): SavedQuery!
	# Delete the saved query with the given ID in the configuration.
	deleteSavedQuery(id: ID!): EmptyResponse
}

# Input to ConfigurationMutation.updateConfiguration. If multiple fields are specified,
# then their respective operations are performed sequentially in the order in which the
# fields appear in this type.
input UpdateConfigurationInput {
	# The name of the property to update.
	#
	# Inserting into an existing array is not yet supported.
	property: String!
	# The new JSON-encoded value to insert. If the field's value is null, the property is
	# removed. (This is different from the field's value being the string "null".)
	value: JSONString
}

# The payload for ConfigurationMutation.updateConfiguration.
type UpdateConfigurationPayload {
	empty: EmptyResponse
}

# The result for Mutation.createUserBySiteAdmin.
type CreateUserBySiteAdminResult {
	# The reset password URL that the new user must visit to sign into their account.
	resetPasswordURL: String!
}

# The result for Mutation.randomizeUserPasswordBySiteAdmin.
type RandomizeUserPasswordBySiteAdminResult {
	# The reset password URL that the user must visit to sign into their account again.
	resetPasswordURL: String!
}

type Query {
	root: Query! @deprecated
	node(id: ID!): Node
	repository(uri: String!): Repository
	phabricatorRepo(uri: String!): PhabricatorRepo
	# A list of all repositories on this site.
	repositories(
		# Returns the first n repositories from the list.
		first: Int
	): RepositoryConnection!
	symbols(id: String!, mode: String!): [Symbol!]!
	currentUser: User
	isUsernameAvailable(username: String!): Boolean!
	configuration: ConfigurationCascade!
	search(query: String = "", scopeQuery: String = ""): Search
	searchScopes: [SearchScope!]!
	# All saved queries configured for the current user, merged from all configurations.
	savedQueries: [SavedQuery!]!
	repoGroups: [RepoGroup!]!
	org(id: ID!): Org! @deprecated(reason: "use Query.node instead")
	sharedItem(ulid: String!): SharedItem
	packages(
		lang: String!
		id: String
		type: String
		name: String
		commit: String
		baseDir: String
		repoURL: String
		version: String
		limit: Int
	): [Package!]!
	dependents(
		lang: String!
		id: String
		type: String
		name: String
		commit: String
		baseDir: String
		repoURL: String
		version: String
		package: String
		limit: Int
	): [Dependency!]!
	users(
		# Returns the first n users from the list.
		first: Int
	): UserConnection!
	# List all organizations.
	orgs(
		# Returns the first n organizations from the list.
		first: Int
	): OrgConnection!
	updateDeploymentConfiguration(email: String!, enableTelemetry: Boolean!): EmptyResponse
	# The current site.
	site: Site!
}

type Search {
	results: SearchResults!
	suggestions(first: Int): [SearchSuggestion!]!
}

union SearchResult = FileMatch | CommitSearchResult

type SearchResults {
	results: [SearchResult!]!
	resultCount: Int!
	approximateResultCount: String!
	limitHit: Boolean!
	# Repositories that are busy cloning onto gitserver.
	cloning: [String!]!
	# Repositories or commits that do not exist.
	missing: [String!]!
	# Repositories or commits which we did not manage to search in time. Trying
	# again usually will work.
	timedout: [String!]!
	# An alert message that should be displayed before any results.
	alert: SearchAlert
}

union SearchSuggestion = Repository | File

type SearchScope {
	name: String!
	value: String!
}

# A search-related alert message.
type SearchAlert {
	title: String!
	description: String
	# "Did you mean: ____" query proposals
	proposedQueries: [SearchQueryDescription!]
}

# A saved search query, defined in configuration.
type SavedQuery {
	# The unique ID of the saved query.
	id: ID!
	# The subject whose configuration this saved query was defined in.
	subject: ConfigurationSubject!
	# The unique key of this saved query (unique only among all other saved
	# queries of the same subject).
	key: String
	# The 0-indexed index of this saved query in the subject's configuration.
	index: Int!
	description: String!
	query: SearchQuery!
}

type SearchQueryDescription {
	description: String
	query: SearchQuery!
}

type SearchQuery {
	query: String!
	scopeQuery: String!
}

# A group of repositories.
type RepoGroup {
	name: String!
	repositories: [String!]!
}

# A diff between two diffable Git objects.
type Diff {
	# The diff's repository.
	repository: Repository!
	# The revision range of the diff.
	range: GitRevisionRange!
}

# A search result that is a Git commit.
type CommitSearchResult {
	# The commit that matched the search query.
	commit: CommitInfo!
	# The ref names of the commit.
	refs: [GitRef!]!
	# The refs by which this commit was reached.
	sourceRefs: [GitRef!]!
	# The matching portion of the commit message, if any.
	messagePreview: HighlightedString
	# The matching portion of the diff, if any.
	diffPreview: HighlightedString
}

# A search result that is a diff between two diffable Git objects.
type DiffSearchResult {
	# The diff that matched the search query.
	diff: Diff!
	# The matching portion of the diff.
	preview: HighlightedString!
}

# A string that has highlights (e.g, query matches).
type HighlightedString {
	# The full contents of the string.
	value: String!
	# Highlighted matches of the query in the preview string.
	highlights: [Highlight!]!
}

# A highlighted region in a string (e.g., matched by a query).
type Highlight {
	# The 1-indexed line number.
	line: Int!
	# The 1-indexed character on the line.
	character: Int!
	# The length of the highlight, in characters (on the same line).
	length: Int!
}

# Represents a shared item (either a shared code comment OR code snippet).
#
# 🚨 SECURITY: Every field here is accessible publicly given a shared item URL.
# Do NOT use any non-primitive graphql type here unless it is also a SharedItem
# type.
type SharedItem {
	# who shared the item.
	author: SharedItemUser!
	public: Boolean!
	thread: SharedItemThread!
	# present only if the shared item was a specific comment.
	comment: SharedItemComment
}

# Like the User type, except with fields that should not be accessible with a
# secret URL removed.
#
# 🚨 SECURITY: Every field here is accessible publicly given a shared item URL.
# Do NOT use any non-primitive graphql type here unless it is also a SharedItem
# type.
type SharedItemUser {
	displayName: String
	username: String!
	avatarURL: String
}

# Like the Thread type, except with fields that should not be accessible with a
# secret URL removed.
#
# 🚨 SECURITY: Every field here is accessible publicly given a shared item URL.
# Do NOT use any non-primitive graphql type here unless it is also a SharedItem
# type.
type SharedItemThread {
	id: Int!
	repo: SharedItemOrgRepo!
	file: String!
	branch: String
	repoRevision: String!
	linesRevision: String!
	title: String!
	startLine: Int!
	endLine: Int!
	startCharacter: Int!
	endCharacter: Int!
	rangeLength: Int!
	createdAt: String!
	archivedAt: String
	author: SharedItemUser!
	lines: SharedItemThreadLines
	comments: [SharedItemComment!]!
}

# Like the OrgRepo type, except with fields that should not be accessible with
# a secret URL removed.
#
# 🚨 SECURITY: Every field here is accessible publicly given a shared item URL.
# Do NOT use any non-primitive graphql type here unless it is also a SharedItem
# type.
type SharedItemOrgRepo {
	id: Int!
	remoteUri: String!
}

# Like the Comment type, except with fields that should not be accessible with a
# secret URL removed.
#
# 🚨 SECURITY: Every field here is accessible publicly given a shared item URL.
# Do NOT use any non-primitive graphql type here unless it is also a SharedItem
# type.
type SharedItemComment {
	id: Int!
	title: String!
	contents: String!
	richHTML: String!
	createdAt: String!
	updatedAt: String!
	author: SharedItemUser!
}

# Exactly the same as the ThreadLines type, except it cannot have sensitive
# fields accidently added.
#
# 🚨 SECURITY: Every field here is accessible publicly given a shared item URL.
# Do NOT use any non-primitive graphql type here unless it is also a SharedItem
# type.
type SharedItemThreadLines {
	htmlBefore(isLightTheme: Boolean!): String!
	html(isLightTheme: Boolean!): String!
	htmlAfter(isLightTheme: Boolean!): String!
	textBefore: String!
	text: String!
	textAfter: String!
	textSelectionRangeStart: Int!
	textSelectionRangeLength: Int!
}

type RefFields {
	refLocation: RefLocation
	uri: URI
}

type URI {
	host: String!
	fragment: String!
	path: String!
	query: String!
	scheme: String!
}

type RefLocation {
	startLineNumber: Int!
	startColumn: Int!
	endLineNumber: Int!
	endColumn: Int!
}

# A list of repositories.
type RepositoryConnection {
	# A list of repositories.
	nodes: [Repository!]!
	# The total count of repositories in the connection.
	totalCount: Int!
}

type Repository implements Node {
	id: ID!
	uri: String!
	description: String!
	language: String!
	fork: Boolean!
	starsCount: Int
	forksCount: Int
	private: Boolean!
	createdAt: String!
	pushedAt: String!
	commit(rev: String!): CommitState!
	revState(rev: String!): RevState!
	latest: CommitState!
	lastIndexedRevOrLatest: CommitState!
	# defaultBranch will not be set if we are cloning the repository
	defaultBranch: String
	branches: [String!]!
	tags: [String!]!
	listTotalRefs: TotalRefList!
	gitCmdRaw(params: [String!]!): String!
	# Link to another sourcegraph instance location where this repository is located.
	redirectURL: String
}

# A Git object ID (SHA-1 hash, 40 hexadecimal characters).
scalar GitObjectID

# A Git ref.
type GitRef {
	# The full ref name (e.g., "refs/heads/mybranch" or "refs/tags/mytag").
	name: String!
	# The display name of the ref. For branches ("refs/heads/foo"), this is the branch
	# name ("foo").
	#
	# As a special case, for GitHub pull request refs of the form refs/pull/NUMBER/head,
	# this is "#NUMBER".
	displayName: String!
	# The prefix of the ref, either "", "refs/", "refs/heads/", "refs/pull/", or
	# "refs/tags/". This prefix is always a prefix of the ref's name.
	prefix: String!
	# The object that the ref points to.
	target: GitObject!
	# The associated repository.
	repository: Repository!
}

# A Git object.
type GitObject {
	oid: GitObjectID!
}

# A Git revspec expression that (possibly) evaluates to a Git revision.
type GitRevSpecExpr {
	expr: String!
}

# A Git revspec.
union GitRevSpec = GitRef | GitRevSpecExpr | GitObject

# A Git revision range of the form "base..head" or "base...head". Other revision
# range formats are not supported.
type GitRevisionRange {
	# The Git revision range expression of the form "base..head" or "base...head".
	expr: String!
	# The base (left-hand side) of the range.
	base: GitRevSpec!
	# The base's revspec as an expression.
	baseRevSpec: GitRevSpecExpr!
	# The head (right-hand side) of the range.
	head: GitRevSpec!
	# The head's revspec as an expression.
	headRevSpec: GitRevSpecExpr!
	# The merge-base of the base and head revisions, if this is a "base...head"
	# revision range. If this is a "base..head" revision range, then this field is null.
	mergeBase: GitObject
}

type PhabricatorRepo {
	# the canonical repo path, like 'github.com/gorilla/mux'
	uri: String!
	# the unique Phabricator identifier for the repo, like 'MUX'
	callsign: String!
	# the URL to the phabricator instance, e.g. http://phabricator.sgdev.org
	url: String!
}

type TotalRefList {
	repositories: [Repository!]!
	total: Int!
}

type Symbol {
	repository: Repository!
	path: String!
	line: Int!
	character: Int!
}

type CommitState {
	commit: Commit
	cloneInProgress: Boolean!
}

type RevState {
	commit: Commit
	cloneInProgress: Boolean!
}

type Commit implements Node {
	id: ID!
	sha1: String!
	tree(path: String = "", recursive: Boolean = false): Tree
	file(path: String!): File
	languages: [String!]!
}

type CommitInfo {
	repository: Repository!
	oid: GitObjectID!
	abbreviatedOID: String!
	rev: String!
	author: Signature!
	committer: Signature
	message: String!
}

type Signature {
	person: Person
	date: String!
}

type Person {
	name: String!
	email: String!
	# The name if set; otherwise the email username.
	displayName: String!
	gravatarHash: String!
	avatarURL: String!
}

type Tree {
	directories: [Directory]!
	files: [File]!
	# Consists of directories plus files.
	entries: [TreeEntry!]!
}

# A file, directory, or other tree entry.
interface TreeEntry {
	name: String!
	isDirectory: Boolean!
	repository: Repository!
	commits: [CommitInfo!]!
	lastCommit: CommitInfo!
}

type Directory implements TreeEntry {
	name: String!
	isDirectory: Boolean!
	repository: Repository!
	commits: [CommitInfo!]!
	lastCommit: CommitInfo!
	tree: Tree!
}

type HighlightedFile {
	aborted: Boolean!
	html: String!
}

type File implements TreeEntry {
	name: String!
	content: String!

	# The file rendered as rich HTML, or an empty string if it is not a supported
	# rich file type.
	#
	# This HTML string is already escaped and thus is always safe to render.
	richHTML: String!

	repository: Repository!
	binary: Boolean!
	isDirectory: Boolean!
	commit: Commit!
	highlight(disableTimeout: Boolean!, isLightTheme: Boolean!): HighlightedFile!
	blame(startLine: Int!, endLine: Int!): [Hunk!]!
	commits: [CommitInfo!]!
	lastCommit: CommitInfo!
	dependencyReferences(Language: String!, Line: Int!, Character: Int!): DependencyReferences!
	blameRaw(startLine: Int!, endLine: Int!): String!
}

type FileMatch {
	resource: String!
	lineMatches: [LineMatch!]!
	limitHit: Boolean!
}

type LineMatch {
	preview: String!
	lineNumber: Int!
	offsetAndLengths: [[Int!]!]!
	limitHit: Boolean!
}

type DependencyReferences {
	dependencyReferenceData: DependencyReferencesData!
	repoData: RepoDataMap!
}

type RepoDataMap {
	repos: [Repository!]!
	repoIds: [Int!]!
}

type DependencyReferencesData {
	references: [DependencyReference!]!
	location: DepLocation!
}

type DependencyReference {
	dependencyData: String!
	repoId: Int!
	hints: String!
}

type DepLocation {
	location: String!
	symbol: String!
}

type Hunk {
	startLine: Int!
	endLine: Int!
	startByte: Int!
	endByte: Int!
	rev: String!
	author: Signature
	message: String!
}

# A list of users.
type UserConnection {
	# A list of users.
	nodes: [User!]!
	# The total count of users in the connection.
	totalCount: Int!
}

type User implements Node, ConfigurationSubject {
	# The unique ID for the user.
	id: ID!
	authID: String!
	auth0ID: String! @deprecated(reason: "use authID instead")
	sourcegraphID: Int!
	email: String!
	displayName: String
	username: String!
	avatarURL: String
	createdAt: String!
	updatedAt: String
	verified: Boolean!
	# Whether the user is a site admin.
	siteAdmin: Boolean!
	# The latest settings for the user.
	latestSettings: Settings
	orgs: [Org!]!
	orgMemberships: [OrgMember!]!
	tags: [UserTag!]!
	activity: UserActivity!
}

# A list of organizations.
type OrgConnection {
	# A list of organizations.
	nodes: [Org!]!
	# The total count of organizations in the connection.
	totalCount: Int!
}

type Org implements Node, ConfigurationSubject {
	id: ID!
	orgID: Int!
	name: String!
	displayName: String
	slackWebhookURL: String
	# The date when the organization was created, in RFC 3339 format.
	createdAt: String!
	members: [OrgMember!]!
	latestSettings: Settings
	repos: [OrgRepo!]!
	repo(canonicalRemoteID: String!): OrgRepo
	threads(
		# TODO(nick): remove repoCanonicalRemoteID
		repoCanonicalRemoteID: String
		canonicalRemoteIDs: [String!]
		branch: String
		file: String
		limit: Int
	): ThreadConnection!
	tags: [OrgTag!]!
}

type OrgMember {
	id: Int!
	org: Org!
	user: User!
	createdAt: String!
	updatedAt: String!
}

type InviteUserResult {
	# The URL that the invited user can visit to accept the invitation.
	acceptInviteURL: String!
}

type OrgInviteStatus {
	emailVerified: Boolean!
}

type OrgRepo {
	id: Int!
	org: Org!
	canonicalRemoteID: String!
	createdAt: String!
	updatedAt: String!
	threads(file: String, branch: String, limit: Int): ThreadConnection!
}

# A list of threads.
type ThreadConnection {
	# A list of threads.
	nodes: [Thread!]!
	# The total count of threads in the connection.
	totalCount: Int!
}

# A site is an installation of Sourcegraph that consists of one or more
# servers that share the same configuration and database.
#
# The site is a singleton; the API only ever returns the single global site.
type Site implements ConfigurationSubject {
	# The site's ID.
	id: ID!
	# The site's configuration. Only visible to site admins.
	configuration: SiteConfiguration!
	# The site's latest site-wide settings (which are the lowest-precedence
	# in the configuration cascade for a user).
	latestSettings: Settings
	# Whether the viewer can reload the site (with the reloadSite mutation).
	canReloadSite: Boolean!
	# List all threads.
	threads(
		# Returns the first n threads from the list.
		first: Int
	): ThreadConnection!
}

# The configuration for a site.
type SiteConfiguration {
	# The effective configuration JSON. This will lag behind the pendingContents
	# if the site configuration was updated but the server has not yet restarted.
	effectiveContents: String!
	# The pending configuration JSON, which will become effective after the next
	# server restart. This is set if the site configuration has been updated since
	# the server started.
	pendingContents: String
	# Validation errors on the configuration JSON (pendingContents if it exists, otherwise
	# effectiveContents). These are different from the JSON Schema validation errors;
	# they are errors from validations that can't be expressed in the JSON Schema.
	extraValidationErrors: [String!]!
	# Whether the viewer can update the site configuration (using the
	# updateSiteConfiguration mutation).
	canUpdate: Boolean!
	# The source of the configuration as a human-readable description,
	# referring to either the on-disk file path or the SOURCEGRAPH_CONFIG
	# env var.
	source: String!
}

# ConfigurationSubject is something that can have configuration.
interface ConfigurationSubject {
	id: ID!
	latestSettings: Settings
}

# The configurations for all of the relevant configuration subjects, plus the merged
# configuration.
type ConfigurationCascade {
	# The default settings, which are applied first and the lowest priority behind
	# all configuration subjects' settings.
	defaults: Configuration
	# The configurations for all of the subjects that are applied for the currently
	# authenticated user. For example, a user in 2 orgs would have the following
	# configuration subjects: org 1, org 2, and the user.
	subjects: [ConfigurationSubject!]!
	# The effective configuration, merged from all of the subjects.
	merged: Configuration!
}

# Settings is a version of a configuration settings file.
type Settings {
	id: Int!
	configuration: Configuration!
	# The subject that these settings are for.
	subject: ConfigurationSubject!
	author: User!
	createdAt: String!
	contents: String! @deprecated(reason: "use configuration.contents instead")
}

# Configuration contains settings from (possibly) multiple settings files.
type Configuration {
	# The raw JSON contents, encoded as a string.
	contents: String!
	# Error and warning messages about the configuration.
	messages: [String!]!
}

type Thread {
	id: Int!
	repo: OrgRepo!
	file: String! @deprecated(reason: "use repoRevisionPath (or linesRevisionPath)")

	# The relative path of the resource in the repository at repoRevision.
	repoRevisionPath: String!

	# The relative path of the resource in the repository at linesRevision.
	linesRevisionPath: String!

	branch: String
	# The commit ID of the repository at the time the thread was created.
	repoRevision: String!
	# The commit ID from Git blame, at the time the thread was created.
	#
	# The selection may be multiple lines, and the commit id is the
	# topologically most recent commit of the blame commit ids for the selected
	# lines.
	#
	# For example, if you have a selection of lines that have blame revisions
	# (a, c, e, f), and assuming a history like::
	#
	# 	a <- b <- c <- d <- e <- f <- g <- h <- HEAD
	#
	# Then lines_revision would be f, because all other blame revisions a, c, e
	# are reachable from f.
	#
	# Or in lay terms: "What is the oldest revision that I could checkout and
	# still see the exact lines of code that I selected?".
	linesRevision: String!
	title: String!
	startLine: Int!
	endLine: Int!
	startCharacter: Int!
	endCharacter: Int!
	rangeLength: Int!
	createdAt: String!
	archivedAt: String
	author: User!
	lines: ThreadLines
	comments: [Comment!]!
}

type Comment {
	id: Int!
	title: String!
	contents: String!

	# The file rendered as rich HTML, or an empty string if it is not a supported
	# rich file type.
	#
	# This HTML string is already escaped and thus is always safe to render.
	richHTML: String!

	createdAt: String!
	updatedAt: String!
	author: User!
}

type Package {
	lang: String!
	repo: Repository
	# The following fields are properties of build package configuration as returned by the workspace/xpackages LSP endpoint.
	id: String
	type: String
	name: String
	commit: String
	baseDir: String
	repoURL: String
	version: String
}

type Dependency {
	repo: Repository
	# The following fields are properties of build package configuration as returned by the workspace/xpackages LSP endpoint.
	name: String
	repoURL: String
	depth: Int
	vendor: Boolean
	package: String
	absolute: String
	type: String
	commit: String
	version: String
	id: String
	package: String
}

type UserTag {
	id: Int!
	name: String!
}

type OrgTag {
	id: Int!
	name: String!
}

type UserActivity {
	id: Int!
	searchQueries: Int!
	pageViews: Int!
	createdAt: String!
	updatedAt: String!
}

enum UserEvent {
	PAGEVIEW
	SEARCHQUERY
}

type DeploymentConfiguration {
	email: String
	telemetryEnabled: Boolean
	appID: String
}
`
