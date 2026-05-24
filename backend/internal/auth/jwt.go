package auth

import(
	"time"
	"github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

type Service struct{
	secret string
	expiry time.Duration
}

func NewService(secret string, expiry time.Duration)*Service{
	return &Service{
		secret: secret,
		expiry: expiry,
	}
}


func(s *Service) Generate(userID uuid.UUID,role string )(string,error){
	claims :=jwt.MapClaims{
		"sub": userID.String(),
		"role": role,
		"exp": time.Now().Add(s.expiry).Unix(),
		"iat": time.Now().Unix(),
	}
	token:=jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	signedToken,err:=token.SignedString([]byte(s.secret))
	
	if err!=nil{
		return "",err
	}
	return signedToken,nil
}