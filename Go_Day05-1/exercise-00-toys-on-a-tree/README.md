В данном коде представлены несколько функций, работающих с бинарным деревом. Вот разъяснение каждой из них:

### 1. Структура `TreeNode`

```go
type TreeNode struct {
    HasToy bool
    Left   *TreeNode
    Right  *TreeNode
}
```

- **Описание**: Эта структура представляет узел бинарного дерева. Каждый узел содержит:
  - `HasToy` (bool): Указывает, есть ли игрушка в узле.
  - `Left` (указатель на `TreeNode`): Указатель на левый дочерний узел.
  - `Right` (указатель на `TreeNode`): Указатель на правый дочерний узел.

### 2. Функция `areToysBalanced`

```go
func areToysBalanced(root *TreeNode) bool {
    if root == nil {
        return true
    }
    return countToys(root.Left) == countToys(root.Right)
}
```

- **Цель**: Проверить, сбалансировано ли количество игрушек между левым и правым поддеревьями корневого узла.
- **Работа**:
  - Если `root` равен `nil`, возвращает `true`, потому что пустое дерево считается сбалансированным.
  - Использует вспомогательную функцию `countToys` для подсчета количества игрушек в левом и правом поддеревьях. Если количество игрушек в обоих поддеревьях равно, дерево сбалансировано, и функция возвращает `true`; иначе — `false`.

### 3. Функция `countToys`

```go
func countToys(node *TreeNode) int {
    if node == nil {
        return 0
    }
    count := 0
    if node.HasToy {
        count = 1
    }
    return count + countToys(node.Left) + countToys(node.Right)
}
```

- **Цель**: Подсчитать общее количество игрушек в поддереве, начиная с указанного узла.
- **Работа**:
  - Если `node` равен `nil`, возвращает `0`.
  - Если узел содержит игрушку (`node.HasToy` равно `true`), увеличивает счетчик на 1.
  - Рекурсивно вызывает `countToys` для левого и правого дочерних узлов и суммирует их количество с текущим значением `count`.

### 4. Функция `drawTree`

```go
func drawTree(node *TreeNode, prefix string, isLast bool) {
    if node != nil {
        // Определяем символ для текущего узла
        var nodeSymbol string
        if node.HasToy {
            nodeSymbol = "1"
        } else {
            nodeSymbol = "0"
        }

        // Отрисовка текущего узла
        if isLast {
            fmt.Println(prefix + "└── " + nodeSymbol)
            prefix += "    "
        } else {
            fmt.Println(prefix + "├── " + nodeSymbol)
            prefix += "│   "
        }

        // Рекурсивно отрисовываем поддеревья
        if node.Left != nil || node.Right != nil {
            drawTree(node.Left, prefix, node.Right == nil)
            drawTree(node.Right, prefix, true)
        }
    }
}
```

- **Цель**: Визуально отобразить бинарное дерево в виде текстового графа.
- **Работа**:
  - Использует префикс для отрисовки уровня и структуры дерева.
  - В зависимости от того, является ли текущий узел последним на этом уровне (`isLast`), используется символ `└──` для последнего узла и `├──` для остальных.
  - Рекурсивно вызывает саму себя для левого и правого дочерних узлов, обновляя префикс для правильного отображения дерева.

### 5. Функция `promptForBool`

```go
func promptForBool(prompt string) bool {
    var input string
    for {
        fmt.Printf("%s (true/false): ", prompt)
        fmt.Scanln(&input)
        if input == "true" {
            return true
        } else if input == "false" {
            return false
        } else {
            fmt.Println("Invalid input. Please enter 'true' or 'false'.")
        }
    }
}
```

- **Цель**: Запросить у пользователя ввод значения типа `bool`.
- **Работа**:
  - Запрашивает ввод от пользователя до тех пор, пока не будет введено корректное значение (`true` или `false`).

### 6. Функция `promptForSubtree`

```go
func promptForSubtree(prompt string) string {
    var input string
    for {
        fmt.Printf("%s (left/right/both/none): ", prompt)
        fmt.Scanln(&input)
        if input == "left" || input == "right" || input == "both" || input == "none" {
            return input
        } else {
            fmt.Println("Invalid input. Please enter 'left', 'right', 'both', or 'none'.")
        }
    }
}
```

- **Цель**: Запросить у пользователя, какие поддеревья должны быть созданы для узла.
- **Работа**:
  - Запрашивает у пользователя ввод, который может быть `left`, `right`, `both`, или `none`, до тех пор, пока не будет введено корректное значение.

### 7. Функция `main`

```go
func main() {
    // Запрашиваем у пользователя значения для каждого узла
    rootToy := promptForBool("Does the root node have a toy?")
    leftToy := promptForBool("Does the left child node have a toy?")
    rightToy := promptForBool("Does the right child node have a toy?")

    // Проверяем, какие поддеревья создать для левого дочернего узла
    var leftSubtree *TreeNode
    leftSubtreeChoice := promptForSubtree("Which subtree does the left child node have?")
    switch leftSubtreeChoice {
    case "left":
        leftLeftToy := promptForBool("  Does the left subtree, left child node have a toy?")
        leftSubtree = &TreeNode{
            HasToy: leftToy,
            Left:   &TreeNode{HasToy: leftLeftToy},
        }
    case "right":
        leftRightToy := promptForBool("  Does the left subtree, right child node have a toy?")
        leftSubtree = &TreeNode{
            HasToy: leftToy,
            Right:  &TreeNode{HasToy: leftRightToy},
        }
    case "both":
        leftLeftToy := promptForBool("  Does the left subtree, left child node have a toy?")
        leftRightToy := promptForBool("  Does the left subtree, right child node have a toy?")
        leftSubtree = &TreeNode{
            HasToy: leftToy,
            Left:   &TreeNode{HasToy: leftLeftToy},
            Right:  &TreeNode{HasToy: leftRightToy},
        }
    default:
        leftSubtree = &TreeNode{HasToy: leftToy}
    }

    // Проверяем, какие поддеревья создать для правого дочернего узла
    var rightSubtree *TreeNode
    rightSubtreeChoice := promptForSubtree("Which subtree does the right child node have?")
    switch rightSubtreeChoice {
    case "left":
        rightLeftToy := promptForBool("  Does the right subtree, left child node have a toy?")
        rightSubtree = &TreeNode{
            HasToy: rightToy,
            Left:   &TreeNode{HasToy: rightLeftToy},
        }
    case "right":
        rightRightToy := promptForBool("  Does the right subtree, right child node have a toy?")
        rightSubtree = &TreeNode{
            HasToy: rightToy,
            Right:  &TreeNode{HasToy: rightRightToy},
        }
    case "both":
        rightLeftToy := promptForBool("  Does the right subtree, left child node have a toy?")
        rightRightToy := promptForBool("  Does the right subtree, right child node have a toy?")
        rightSubtree = &TreeNode{
            HasToy: rightToy,
            Left:   &TreeNode{HasToy: rightLeftToy},
            Right:  &TreeNode{HasToy: rightRightToy},
        }
    default:
        rightSubtree = &TreeNode{HasToy: rightToy}
    }

    // Создаем дерево на основе введенных данных
    root := &TreeNode{
        HasToy: rootToy,
        Left:   leftSubtree,
        Right:  rightSubtree,
    }

    // Отрисовываем дерево в терминале
    fmt.Println("Binary Tree:")
    drawTree(root, "", true)

    // Проверяем, уравновешено ли дерево
    if areToysBalanced(root) {
        fmt.Println("The tree is balanced - true.")
    } else {
        fmt.Println("The tree is not balanced - false.")
    }
}
```

- **Цель**: Построить бинарное дерево на основе ввода пользователя, отобразить его и проверить, сбалансировано ли дерево по количеству игрушек между левым и правым под

деревьями.
- **Работа**:
  - Запрашивает у пользователя информацию о каждом узле дерева (есть ли игрушка, какие поддеревья создавать).
  - Создает бинарное дерево на основе введенных данных.
  - Отображает дерево в текстовом виде с помощью функции `drawTree`.
  - Проверяет, сбалансировано ли дерево по количеству игрушек, используя функцию `areToysBalanced`.

Таким образом, вы можете создавать и визуализировать бинарное дерево, а также проверять сбалансированность по количеству игрушек в левой и правой ветвях.