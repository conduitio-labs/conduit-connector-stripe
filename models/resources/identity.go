// Copyright © 2022 Meroxa, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package resources

const (
	IdentityVerificationSessionResource           = "identity.verification_session"
	IdentityVerificationSessionsList              = "identity/verification_sessions"
	IdentityVerificationSessionCanceledEvent      = "identity.verification_session.canceled"
	IdentityVerificationSessionCreatedEvent       = "identity.verification_session.created"
	IdentityVerificationSessionProcessingEvent    = "identity.verification_session.processing"
	IdentityVerificationSessionRedactedEvent      = "identity.verification_session.redacted"
	IdentityVerificationSessionRequiresInputEvent = "identity.verification_session.requires_input"
	IdentityVerificationSessionVerifiedEvent      = "identity.verification_session.verified"

	IdentityVerificationReportResource = "identity.verification_report"
	IdentityVerificationReportsList    = "identity/verification_reports"
)

var (
	IdentityVerificationSessionEvents = []string{
		IdentityVerificationSessionCanceledEvent,
		IdentityVerificationSessionCreatedEvent,
		IdentityVerificationSessionProcessingEvent,
		IdentityVerificationSessionRedactedEvent,
		IdentityVerificationSessionRequiresInputEvent,
		IdentityVerificationSessionVerifiedEvent,
	}
)
