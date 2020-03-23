package sales

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

const  (
	createTable = `
create table if not exists sales (
    id BIGSERIAL primary key,
    user_id integer not null,
    product_id integer not null,
    price integer not null check ( price>=0 ),
    quantity integer not null check ( quantity>0 ),
    purchase_date date default current_date,
    removed BOOLEAN DEFAULT FALSE
);
`
	addSales = `INSERT INTO sales(user_id, product_id, price, quantity)
VALUES ($1, $2, $3, $4);`

	salesList = `select id, user_id, product_id, price, quantity, purchase_date from sales where removed=false;`

	removeId = `update sales set removed = true where id = $1`
	buyId = `select id, user_id, product_id, price, quantity, purchase_date from sales where id=$1`
	updateData = `update sales set purchase_date = CURRENT_TIMESTAMP where id = $1`
	updateQuantity = `update sales set quantity = $2 where id = $1`

	updatePrice =`update sales set price = $2 where id = $1`
	updateSales =`update sales set product_id = $2 where id = $1`
	)
type Service struct {
	pool *pgxpool.Pool
}

func (s *Service) Start() {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		panic(errors.New("can't create database"))
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), createTable)
	if err != nil {
		panic(errors.New("can't create database"))
	}

}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

type Purchase struct {
	ID         int64     `json:"id"`
	User_id    int64     `json:"user_id"`
	Product_id int64     `json:"product_id"`
	Price      int       `json:"price"`
	Quantity   int       `json:"quantity"`
	Date       time.Time `json:"date"`
}

func (s *Service) AddNewPurchase(ctx context.Context, prod Purchase) (err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()
	_, err = conn.Exec(ctx, addSales, prod.User_id, prod.Product_id, prod.Price, prod.Quantity)
	if err != nil {
		return
	}
	return nil
}
func (s *Service) BuyByID(ctx context.Context, id int64) (prod Purchase, err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return Purchase{}, errors.New("can't connect to database!")
	}
	defer conn.Release()
	err = conn.QueryRow(ctx,
		buyId,
		id).Scan(&prod.ID, &prod.User_id, &prod.Product_id, &prod.Price, &prod.Quantity, &prod.Date)
	if err != nil {
		return Purchase{}, errors.New(fmt.Sprintf("can't remove from database burger (id: %d)!", id))
	}
	return
}

func (s *Service) BuyingList(ctx context.Context) (list []Purchase, err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(ctx,
		salesList)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := Purchase{}
		err := rows.Scan(&item.ID, &item.User_id, &item.Product_id, &item.Price, &item.Quantity, &item.Date)
		if err != nil {
			return nil, errors.New("can't scan row from rows")
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.New("rows error!")
	}
	return
}
func (s *Service) UpdateBought(ctx context.Context, id int64, pur Purchase) (err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		err = errors.New("can't connect to database!")
		return
	}
	defer conn.Release()
	begin, err := conn.Begin(ctx)
	if err != nil {
		err = errors.New("can't connect to database!")
		return
	}
	defer func() {
		if err != nil {
			err2 := begin.Rollback(ctx)
			if err2 != nil {
				log.Printf("can't rollback err %v", err2)
			}
			return
		}
		err2 := begin.Commit(ctx)
		if err2 != nil {
			log.Printf("can't commit err %v", err2)
		}
	}()
	_, err = begin.Exec(ctx, updateData, id)
	if err != nil {
		return
	}
	if pur.Quantity != -1 {
		_, err = begin.Exec(ctx, updateQuantity, id, pur.Quantity)
		if err != nil {
			return
		}
	}
	if pur.Price != -1 {
		_, err = begin.Exec(ctx, updatePrice, id, pur.Price)
		if err != nil {
			return
		}
	}
	if pur.Product_id != -1 {
		_, err = begin.Exec(ctx, updateSales, id, pur.Product_id)
		if err != nil {
			return
		}
	}
	return
}


func (s *Service) RemoveByID(ctx context.Context, id int64) (err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return errors.New("can't connect to database!")
	}
	defer conn.Release()
	_, err = conn.Exec(ctx, removeId, id)
	if err != nil {
		return errors.New(fmt.Sprintf("can't remove from database sales (id: %d)!", id))
	}
	return nil
}


func (s *Service) BuyByUserID(ctx context.Context, id int64) (list []Purchase, err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		err = errors.New("can't connect to database!")
		return
	}
	defer conn.Release()
	rows, err := conn.Query(ctx,
		`select id, user_id, product_id, price, quantity, purchase_date from sales where user_id=$1`,
		id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := Purchase{}
		err := rows.Scan(&item.ID, &item.User_id, &item.Product_id, &item.Price, &item.Quantity, &item.Date)
		if err != nil {
			return nil, errors.New("can't scan row from rows")
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.New("rows error!")
	}

	return
}

