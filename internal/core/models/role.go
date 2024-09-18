package models

type Role uint8

const (
    InvalidRole Role = 0
    UserRole    Role = 1
    AdminRole   Role = 10
)

func (r Role) String() string {
    switch r {
    case UserRole:
        return "user"
    case AdminRole:
        return "admin"
    default:
        return "unknown"
    }
}

func ParseRole(s string) Role {
    switch s {
    case "user":
        return UserRole
    case "admin":
        return AdminRole
    default:
        return InvalidRole
    }
}
