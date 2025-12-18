package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

// BankClient интерфейс банковского клиента
type BankClient interface {
	// Deposit зачисляет указанную сумму на счет клиента
	Deposit(amount int)
	
	// Withdrawal снимает указанную сумму со счета клиента.
	// Возвращает ошибку, если баланс клиента меньше суммы снятия
	Withdrawal(amount int) error
	
	// Balance возвращает баланс клиента
	Balance() int
}

// bankClient реализация банковского клиента
type bankClient struct {
	balance int
	mutex   sync.Mutex // защищает доступ к балансу
}

// NewBankClient создает новый экземпляр банковского клиента
func NewBankClient() BankClient {
	return &bankClient{
		balance: 0,
	}
}

// Deposit зачисляет указанную сумму на счет клиента
func (bc *bankClient) Deposit(amount int) {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	bc.balance += amount
	fmt.Printf("Зачислено %d. Текущий баланс: %d\n", amount, bc.balance)
}

// Withdrawal снимает указанную сумму со счета клиента
func (bc *bankClient) Withdrawal(amount int) error {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	if bc.balance < amount {
		return fmt.Errorf("недостаточно средств. Текущий баланс: %d, запрашиваемая сумма: %d", bc.balance, amount)
	}
	bc.balance -= amount
	fmt.Printf("Снято %d. Текущий баланс: %d\n", amount, bc.balance)
	return nil
}

// Balance возвращает баланс клиента
func (bc *bankClient) Balance() int {
	bc.mutex.Lock()
	defer bc.mutex.Unlock()
	return bc.balance
}

func main() {
	client := NewBankClient()
	var wg sync.WaitGroup
	stopChan := make(chan struct{})

	// Запуск горутин для зачисления средств (10 штук)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-stopChan:
					return
				default:
					// Случайная сумма от 1 до 10
					amount := rand.Intn(10) + 1
					client.Deposit(amount)
					
					// Ожидание от 0.5 до 1 секунды
					time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
				}
			}
		}(i)
	}

	// Запуск горутин для снятия средств (5 штук)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for {
				select {
				case <-stopChan:
					return
				default:
					// Случайная сумма от 1 до 5
					amount := rand.Intn(5) + 1
					err := client.Withdrawal(amount)
					if err != nil {
						fmt.Println(err.Error())
					}
					
					// Ожидание от 0.5 до 1 секунды
					time.Sleep(time.Duration(rand.Intn(500)+500) * time.Millisecond)
				}
			}
		}(i)
	}

	// Запускаем таймер на 15 секунд для автоматических операций
	fmt.Println("Запуск автоматических операций на 15 секунд...")
	time.AfterFunc(15*time.Second, func() {
		close(stopChan)
	})

	// Ждем завершения автоматических операций
	wg.Wait()
	fmt.Println("Автоматические операции завершены. Теперь можно вводить команды.")

	// Обработка пользовательского ввода
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Доступные команды: balance, deposit, withdrawal, exit")
	
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		
		input := strings.TrimSpace(scanner.Text())
		parts := strings.Split(input, " ")
		command := strings.ToLower(parts[0])
		
		switch command {
		case "balance":
			balance := client.Balance()
			fmt.Printf("Текущий баланс: %d\n", balance)
			
		case "deposit":
			if len(parts) != 2 {
				fmt.Println("Неверный формат команды. Использование: deposit <сумма>")
				continue
			}
			amount, err := strconv.Atoi(parts[1])
			if err != nil || amount <= 0 {
				fmt.Println("Неверная сумма. Введите положительное целое число.")
				continue
			}
			client.Deposit(amount)
			
		case "withdrawal":
			if len(parts) != 2 {
				fmt.Println("Неверный формат команды. Использование: withdrawal <сумма>")
				continue
			}
			amount, err := strconv.Atoi(parts[1])
			if err != nil || amount <= 0 {
				fmt.Println("Неверная сумма. Введите положительное целое число.")
				continue
			}
			err = client.Withdrawal(amount)
			if err != nil {
				fmt.Println(err.Error())
			}
			
		case "exit":
			fmt.Println("Завершение работы приложения...")
			return
			
		default:
			fmt.Println("Неподдерживаемая команда. Доступные команды: balance, deposit, withdrawal, exit")
		}
	}
}