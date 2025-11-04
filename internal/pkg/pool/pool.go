package pool

type Resetter interface {
	Reset()
}

// Pool - пул объектов для переиспользования
type Pool[T Resetter] struct {
	pool []T
	new  func() T
}

// New создает новый пул объектов
func New[T Resetter](newFunc func() T) *Pool[T] {
	return &Pool[T]{
		pool: make([]T, 0),
		new:  newFunc,
	}
}

// Get возвращает объект из пула
func (p *Pool[T]) Get() T {
	if len(p.pool) == 0 {
		return p.new()
	}

	obj := p.pool[len(p.pool)-1]
	p.pool = p.pool[:len(p.pool)-1]
	return obj
}

// Put помещает объект обратно в пул
func (p *Pool[T]) Put(obj T) {
	obj.Reset()
	p.pool = append(p.pool, obj)
}
