package book

type Service interface {
	FindAll() ([]*Book, error)
	FindByID(ID int) (Book, error)
	Create(bookRequest BookRequest) (Book, error)
	Update(ID int, bookRequest BookRequest) (Book, error)
	Delete(ID int) (Book, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) FindAll() ([]*Book, error) {
	return s.repository.FindAll()
}

func (s *service) FindByID(ID int) (Book, error) {
	return s.repository.FindByID(ID)
}

func (s *service) Create(bookRequest BookRequest) (Book, error) {
	price, _ := bookRequest.Price.Int64()
	rating, _ := bookRequest.Rating.Int64()

	book := Book{
		Title:       bookRequest.Title,
		Price:       int(price),
		Rating:      int(rating),
		Description: bookRequest.Description,
	}

	newBook, err := s.repository.Create(book)

	return newBook, err
}

func (s *service) Update(ID int, bookRequest BookRequest) (Book, error) {
	b, findBookErr := s.repository.FindByID(ID)

	if findBookErr != nil {
		return b, findBookErr
	}

	price, _ := bookRequest.Price.Int64()
	rating, _ := bookRequest.Rating.Int64()

	b.Title = bookRequest.Title
	b.Price = int(price)
	b.Description = bookRequest.Description
	b.Rating = int(rating)

	newBook, err := s.repository.Update(b)

	return newBook, err
}

func (s *service) Delete(ID int) (Book, error) {
	b, findBookErr := s.repository.FindByID(ID)

	if findBookErr != nil {
		return b, findBookErr
	}

	book, err := s.repository.Delete(b)

	return book, err
}
