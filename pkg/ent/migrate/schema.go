// Code generated by entc, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// PermissionsColumns holds the columns for the "permissions" table.
	PermissionsColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true, Size: 35},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "delete_time", Type: field.TypeTime, Nullable: true},
		{Name: "name", Type: field.TypeString},
		{Name: "slug", Type: field.TypeString},
		{Name: "status", Type: field.TypeInt32, Default: 1},
	}
	// PermissionsTable holds the schema information for the "permissions" table.
	PermissionsTable = &schema.Table{
		Name:       "permissions",
		Columns:    PermissionsColumns,
		PrimaryKey: []*schema.Column{PermissionsColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "permission_create_time_update_time_delete_time",
				Unique:  false,
				Columns: []*schema.Column{PermissionsColumns[1], PermissionsColumns[2], PermissionsColumns[3]},
			},
			{
				Name:    "permission_name_slug_status",
				Unique:  false,
				Columns: []*schema.Column{PermissionsColumns[4], PermissionsColumns[5], PermissionsColumns[6]},
			},
		},
	}
	// RolesColumns holds the columns for the "roles" table.
	RolesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true, Size: 35},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "delete_time", Type: field.TypeTime, Nullable: true},
		{Name: "name", Type: field.TypeString},
		{Name: "slug", Type: field.TypeString},
		{Name: "status", Type: field.TypeInt32, Default: 1},
	}
	// RolesTable holds the schema information for the "roles" table.
	RolesTable = &schema.Table{
		Name:       "roles",
		Columns:    RolesColumns,
		PrimaryKey: []*schema.Column{RolesColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "role_create_time_update_time_delete_time",
				Unique:  false,
				Columns: []*schema.Column{RolesColumns[1], RolesColumns[2], RolesColumns[3]},
			},
			{
				Name:    "role_name_slug_status",
				Unique:  false,
				Columns: []*schema.Column{RolesColumns[4], RolesColumns[5], RolesColumns[6]},
			},
		},
	}
	// UsersColumns holds the columns for the "users" table.
	UsersColumns = []*schema.Column{
		{Name: "id", Type: field.TypeString, Unique: true, Size: 35},
		{Name: "create_time", Type: field.TypeTime},
		{Name: "update_time", Type: field.TypeTime},
		{Name: "delete_time", Type: field.TypeTime, Nullable: true},
		{Name: "name", Type: field.TypeString},
		{Name: "email", Type: field.TypeString},
		{Name: "password", Type: field.TypeString},
		{Name: "status", Type: field.TypeInt32, Default: 1},
	}
	// UsersTable holds the schema information for the "users" table.
	UsersTable = &schema.Table{
		Name:       "users",
		Columns:    UsersColumns,
		PrimaryKey: []*schema.Column{UsersColumns[0]},
		Indexes: []*schema.Index{
			{
				Name:    "user_create_time_update_time_delete_time",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[1], UsersColumns[2], UsersColumns[3]},
			},
			{
				Name:    "user_email_status",
				Unique:  false,
				Columns: []*schema.Column{UsersColumns[5], UsersColumns[7]},
			},
		},
	}
	// PermissionRolesColumns holds the columns for the "permission_roles" table.
	PermissionRolesColumns = []*schema.Column{
		{Name: "permission_id", Type: field.TypeString, Size: 35},
		{Name: "role_id", Type: field.TypeString, Size: 35},
	}
	// PermissionRolesTable holds the schema information for the "permission_roles" table.
	PermissionRolesTable = &schema.Table{
		Name:       "permission_roles",
		Columns:    PermissionRolesColumns,
		PrimaryKey: []*schema.Column{PermissionRolesColumns[0], PermissionRolesColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "permission_roles_permission_id",
				Columns:    []*schema.Column{PermissionRolesColumns[0]},
				RefColumns: []*schema.Column{PermissionsColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "permission_roles_role_id",
				Columns:    []*schema.Column{PermissionRolesColumns[1]},
				RefColumns: []*schema.Column{RolesColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// RoleUsersColumns holds the columns for the "role_users" table.
	RoleUsersColumns = []*schema.Column{
		{Name: "role_id", Type: field.TypeString, Size: 35},
		{Name: "user_id", Type: field.TypeString, Size: 35},
	}
	// RoleUsersTable holds the schema information for the "role_users" table.
	RoleUsersTable = &schema.Table{
		Name:       "role_users",
		Columns:    RoleUsersColumns,
		PrimaryKey: []*schema.Column{RoleUsersColumns[0], RoleUsersColumns[1]},
		ForeignKeys: []*schema.ForeignKey{
			{
				Symbol:     "role_users_role_id",
				Columns:    []*schema.Column{RoleUsersColumns[0]},
				RefColumns: []*schema.Column{RolesColumns[0]},
				OnDelete:   schema.Cascade,
			},
			{
				Symbol:     "role_users_user_id",
				Columns:    []*schema.Column{RoleUsersColumns[1]},
				RefColumns: []*schema.Column{UsersColumns[0]},
				OnDelete:   schema.Cascade,
			},
		},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		PermissionsTable,
		RolesTable,
		UsersTable,
		PermissionRolesTable,
		RoleUsersTable,
	}
)

func init() {
	PermissionRolesTable.ForeignKeys[0].RefTable = PermissionsTable
	PermissionRolesTable.ForeignKeys[1].RefTable = RolesTable
	RoleUsersTable.ForeignKeys[0].RefTable = RolesTable
	RoleUsersTable.ForeignKeys[1].RefTable = UsersTable
}
