package auth

const (
	Email          IdentityType = "email"
	Username       IdentityType = "username"
	BearerIdentity IdentityType = "bearer"
	//PhoneNumber       IdentityType = "phone"
	//APIKey            IdentityType = "apikey"
	//BlockchainAddress IdentityType = "blockchain"

	Password ProofType = "password"
	Token    ProofType = "token"
	//OTP           ProofType = "otp"
	//OAuthToken    ProofType = "oauth_token"
	//APISecret     ProofType = "api_secret"
	//SignedMessage ProofType = "signed_message"
)

type (
	IdentityType string
	ProofType    string

	Authenticator[UserID comparable] interface {
		Authenticate(Payload) (UserID, error)
	}
	AuthenticatorFunc[UserID comparable] func(Payload) (UserID, error)

	Payload struct {
		IdentityType IdentityType
		Identity     string
		Proofs       []Proof
	}
	Proof struct {
		Type  ProofType
		Value string
	}
)

func (fn AuthenticatorFunc[UserID]) Authenticate(p Payload) (UserID, error) {
	return fn(p)
}

// TODO: multi-method authenticator

func New[UserID comparable](
	it IdentityType, pt ProofType,
	check func(identity, proof string) (UserID, error),
) Authenticator[UserID] {
	return AuthenticatorFunc[UserID](func(p Payload) (UserID, error) {
		if p.IdentityType != it || len(p.Proofs) != 1 || p.Proofs[0].Type != pt {
			var zero UserID
			return zero, ErrUnsupportedType
		}
		return check(p.Identity, p.Proofs[0].Value)
	})
}

func SingleUser[UserID comparable](
	username, password string, userID UserID,
) Authenticator[UserID] {
	return New(Username, Password, func(identity, proof string) (UserID, error) {
		if username != identity || password != proof {
			var zero UserID
			return zero, ErrBadCredentials
		}
		return userID, nil
	})
}

func NewPayload(
	it IdentityType, identity string, pt ProofType, proofValue string,
) Payload {
	return Payload{
		IdentityType: it,
		Identity:     identity,
		Proofs:       []Proof{{Type: pt, Value: proofValue}},
	}
}
