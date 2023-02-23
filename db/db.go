package db

import (
	"database/sql"
	"fmt"
	"github.com/Gretamass/kys-backend/provider"
	"github.com/Gretamass/kys-backend/sneaker"
	"github.com/Gretamass/kys-backend/user"
	_ "modernc.org/sqlite"
	"strings"
)

type DB struct {
	db *sql.DB
}

func ConnectDatabase() (*DB, error) {
	db, err := sql.Open("sqlite", "./sqlite.db")
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &DB{
		db: db,
	}, nil
}

// USER methods

func (d *DB) GetUsers() ([]user.User, error) {
	rows, err := d.db.Query("SELECT * FROM users")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	users := make([]user.User, 0)

	for rows.Next() {
		singleUser := user.User{}
		err = rows.Scan(&singleUser.Id, &singleUser.Email, &singleUser.Password, &singleUser.CreatedAt)

		if err != nil {
			return nil, err
		}

		users = append(users, singleUser)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (d *DB) GetUserById(userId int) (user.User, error) {
	row := d.db.QueryRow("SELECT * FROM users WHERE id = ?", userId)

	singleUser := user.User{}
	err := row.Scan(&singleUser.Id, &singleUser.Email, &singleUser.Password, &singleUser.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return user.User{}, fmt.Errorf("no rows found with id %d", userId)
		}
		return user.User{}, err
	}

	return singleUser, nil
}

func (d *DB) AddUser(user user.User) error {
	row, err := d.db.Prepare("INSERT INTO users (email, password) VALUES (?, ?)")

	if err != nil {
		return err
	}

	_, err = row.Exec(user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) UpdateUser(userId int, request user.User) error {
	query := "UPDATE users SET "
	var args []interface{}

	if request.Email != "" {
		query += "email = ?, "
		args = append(args, request.Email)
	}

	if request.Password != "" {
		query += "password = ?, "
		args = append(args, request.Password)
	}

	query = strings.TrimRight(query, ", ")
	query += " WHERE id = ?"
	args = append(args, userId)

	row, err := d.db.Prepare(query)

	if err != nil {
		return err
	}

	_, err = row.Exec(args...)

	if err != nil {
		return err
	}

	return nil
}

func (d *DB) DeleteUser(userId int) error {
	row, err := d.db.Prepare("DELETE FROM users WHERE id = ?")

	if err != nil {
		return err
	}

	result, err := row.Exec(userId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no users found with id %d", userId)
	}

	return nil
}

func (d *DB) LoginUser(user user.User) (bool, error) {
	var userExists bool

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email=? AND password=?)"
	err := d.db.QueryRow(query, user.Email, user.Password).Scan(&userExists)
	if err != nil {
		return false, err
	}

	return userExists, nil
}

// ADMIN methods

func (d *DB) GetAdmins() ([]user.Admin, error) {
	rows, err := d.db.Query("SELECT * FROM admins")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	admins := make([]user.Admin, 0)

	for rows.Next() {
		singleAdmin := user.Admin{}
		err = rows.Scan(&singleAdmin.Id, &singleAdmin.Email, &singleAdmin.Password)

		if err != nil {
			return nil, err
		}

		admins = append(admins, singleAdmin)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return admins, nil
}

func (d *DB) GetAdminById(adminId int) (user.Admin, error) {
	row := d.db.QueryRow("SELECT * FROM admins WHERE id = ?", adminId)

	singleAdmin := user.Admin{}
	err := row.Scan(&singleAdmin.Id, &singleAdmin.Email, &singleAdmin.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			return user.Admin{}, fmt.Errorf("no rows found with id %d", adminId)
		}
		return user.Admin{}, err
	}

	return singleAdmin, nil
}

func (d *DB) AddAdmin(admin user.Admin) error {
	row, err := d.db.Prepare("INSERT INTO admins (email, password) VALUES (?, ?)")

	if err != nil {
		return err
	}

	_, err = row.Exec(admin.Email, admin.Password)
	if err != nil {
		return err
	}

	return nil
}

func (d *DB) UpdateAdmin(adminId int, request user.Admin) error {
	query := "UPDATE admins SET "
	var args []interface{}

	if request.Email != "" {
		query += "email = ?, "
		args = append(args, request.Email)
	}

	if request.Password != "" {
		query += "password = ?, "
		args = append(args, request.Password)
	}

	query = strings.TrimRight(query, ", ")
	query += " WHERE id = ?"
	args = append(args, adminId)

	row, err := d.db.Prepare(query)

	if err != nil {
		return err
	}

	_, err = row.Exec(args...)

	if err != nil {
		return err
	}

	return nil
}

func (d *DB) DeleteAdmin(adminId int) error {
	row, err := d.db.Prepare("DELETE FROM admins WHERE id = ?")

	if err != nil {
		return err
	}

	result, err := row.Exec(adminId)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no admins found with id %d", adminId)
	}

	return nil
}

// ADMIN methods

func (d *DB) GetSneakers() ([]sneaker.Sneaker, error) {
	rows, err := d.db.Query("SELECT * FROM sneakers")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sneakers := make([]sneaker.Sneaker, 0)

	for rows.Next() {
		singleSneaker := sneaker.Sneaker{}
		err = rows.Scan(&singleSneaker.Id, &singleSneaker.Name, &singleSneaker.Model, &singleSneaker.Brand)

		if err != nil {
			return nil, err
		}

		sneakers = append(sneakers, singleSneaker)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return sneakers, nil
}

func (d *DB) GetSneakersInfo() ([]sneaker.SneakerInformation, error) {
	rows, err := d.db.Query("SELECT s.*, si.* FROM sneakers s JOIN sneakers_information si ON s.id = si.sneakerId")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sneakers := make([]sneaker.SneakerInformation, 0)

	for rows.Next() {
		singleSneaker := sneaker.SneakerInformation{}
		err = rows.Scan(&singleSneaker.Id, &singleSneaker.Name, &singleSneaker.Brand, &singleSneaker.Model,
			&singleSneaker.SneakerInformation.SneakerId, &singleSneaker.SneakerInformation.MainInfo,
			&singleSneaker.SneakerInformation.MainImageUrl, &singleSneaker.SneakerInformation.AdditionalInfo)

		if err != nil {
			return nil, err
		}

		sneakers = append(sneakers, singleSneaker)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return sneakers, nil
}

func (d *DB) GetSneakerInfo(sneakerId int) (*sneaker.SneakerInformation, error) {
	query := `
        SELECT s.*, si.*
        FROM sneakers s
        JOIN sneakers_information si ON s.id = si.sneakerId
        WHERE s.id = ?;
    `
	row := d.db.QueryRow(query, sneakerId)

	sneaker := &sneaker.SneakerInformation{}
	err := row.Scan(&sneaker.Id, &sneaker.Name, &sneaker.Brand, &sneaker.Model, &sneaker.SneakerInformation.SneakerId,
		&sneaker.SneakerInformation.MainInfo, &sneaker.SneakerInformation.MainImageUrl, &sneaker.SneakerInformation.AdditionalInfo)
	if err != nil {
		return nil, err
	}

	return sneaker, nil
}

func (d *DB) GetSneakersAvailability() ([]sneaker.SneakerAvailability, error) {
	rows, err := d.db.Query("SELECT s.*, pi.* FROM sneakers s JOIN provider_information pi ON s.id = pi.product_id")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sneakersMap := make(map[int]sneaker.SneakerAvailability)

	for rows.Next() {
		singleSneaker := sneaker.SneakerAvailability{}
		availability := sneaker.Availability{}
		err = rows.Scan(&singleSneaker.Id, &singleSneaker.Name, &singleSneaker.Brand, &singleSneaker.Model,
			&availability.Id, &availability.ProductId, &availability.ProviderId,
			&availability.Price, &availability.Available)

		if err != nil {
			return nil, err
		}

		// Check if we've already added the sneaker to the map
		if existingSneaker, ok := sneakersMap[singleSneaker.Id]; ok {
			// Add the availability to the existing sneaker's array
			existingSneaker.Availability = append(existingSneaker.Availability, availability)
			sneakersMap[singleSneaker.Id] = existingSneaker
		} else {
			// Add the sneaker to the map
			singleSneaker.Availability = []sneaker.Availability{availability}
			sneakersMap[singleSneaker.Id] = singleSneaker
		}
	}

	// Convert the map to an array
	sneakers := make([]sneaker.SneakerAvailability, 0, len(sneakersMap))
	for _, value := range sneakersMap {
		sneakers = append(sneakers, value)
	}

	return sneakers, nil
}

func (d *DB) GetSneakerScrapper(sneakerId int) ([]sneaker.AvailabilityScrappers, error) {
	query := `
        SELECT s.*, avs.*
        FROM sneakers s
        JOIN availability_scrappers avs ON s.id = avs.product_id 
        WHERE s.id = ?;
    `
	rows, err := d.db.Query(query, sneakerId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sneakersMap := make(map[int]sneaker.AvailabilityScrappers)

	for rows.Next() {
		singleSneaker := sneaker.AvailabilityScrappers{}
		scrapper := sneaker.Scrapper{}
		err = rows.Scan(&singleSneaker.Id, &singleSneaker.Name, &singleSneaker.Brand, &singleSneaker.Model,
			&scrapper.Id, &scrapper.ProductId, &scrapper.ProviderId, &scrapper.SearchFor)

		if err != nil {
			return nil, err
		}

		// Check if we've already added the scrapper to the map
		if existingSneaker, ok := sneakersMap[singleSneaker.Id]; ok {
			// Add the scrapper to the existing scrapper's array
			existingSneaker.Scrapper = append(existingSneaker.Scrapper, scrapper)
			sneakersMap[singleSneaker.Id] = existingSneaker
		} else {
			// Add the scrapper to the map
			singleSneaker.Scrapper = []sneaker.Scrapper{scrapper}
			sneakersMap[singleSneaker.Id] = singleSneaker
		}
	}

	// Convert the map to an array
	sneakers := make([]sneaker.AvailabilityScrappers, 0, len(sneakersMap))
	for _, value := range sneakersMap {
		sneakers = append(sneakers, value)
	}

	return sneakers, nil
}

// PROVIDER methods

func (d *DB) GetProviders() ([]provider.ProviderInformation, error) {
	rows, err := d.db.Query("SELECT * FROM product_providers")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	providers := make([]provider.ProviderInformation, 0)

	for rows.Next() {
		singleProvider := provider.ProviderInformation{}
		err = rows.Scan(&singleProvider.Id, &singleProvider.ProviderName)

		if err != nil {
			return nil, err
		}

		providers = append(providers, singleProvider)
	}

	err = rows.Err()

	if err != nil {
		return nil, err
	}

	return providers, nil
}

func (d *DB) GetProviderById(providerId int) (provider.ProviderInformation, error) {
	row := d.db.QueryRow("SELECT * FROM product_providers WHERE id = ?", providerId)

	singleProvider := provider.ProviderInformation{}
	err := row.Scan(&singleProvider.Id, &singleProvider.ProviderName)

	if err != nil {
		if err == sql.ErrNoRows {
			return provider.ProviderInformation{}, fmt.Errorf("no rows found with id %d", providerId)
		}
		return provider.ProviderInformation{}, err
	}

	return singleProvider, nil
}
