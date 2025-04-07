package service

type Place struct {
	name string
}

type Places struct {
	places []Place
}

func NewPlaces() *Places {
	return &Places{
		places: []Place{
			{name: "Школа"},
			{name: "Спортзал"},
			{name: "Кинотеатр"},
			{name: "Музей"},
			{name: "Парк"},
			{name: "Библиотека"},
			{name: "Кафе"},
			{name: "Ресторан"},
			{name: "Театр"},
			{name: "Зоопарк"},
			{name: "Аквапарк"},
			{name: "Торговый центр"},
			{name: "Стадион"},
			{name: "Бассейн"},
			{name: "Выставочный зал"},
			{name: "Концертный зал"},
			{name: "Ботанический сад"},
			{name: "Пляж"},
			{name: "Ночной клуб"},
			{name: "Церковь"},
		},
	}
}

func (p *Places) Get(inx int) Place {
	return p.places[inx]
}
