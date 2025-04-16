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

	Authenticator[User any] interface {
		Authenticate(Payload) (User, error)
	}
	AuthenticatorFunc[User any] func(Payload) (User, error)

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

func (fn AuthenticatorFunc[User]) Authenticate(p Payload) (User, error) {
	return fn(p)
}

// TODO: multi-method authenticator

func New[User any](
	it IdentityType, pt ProofType,
	check func(identity, proof string) (User, error),
) Authenticator[User] {
	return AuthenticatorFunc[User](func(p Payload) (User, error) {
		if p.IdentityType != it || len(p.Proofs) != 1 || p.Proofs[0].Type != pt {
			var zero User
			return zero, ErrUnsupportedType
		}
		return check(p.Identity, p.Proofs[0].Value)
	})
}

func SingleUser[User any](
	username, password string, user User,
) Authenticator[User] {
	return New(Username, Password, func(identity, proof string) (User, error) {
		if username != identity || password != proof {
			var zero User
			return zero, ErrBadCredentials
		}
		return user, nil
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
