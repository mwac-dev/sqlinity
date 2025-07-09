# sqlinity

**SQL ➜ C# migration code generator for Unity and beyond.**  
`sqlinity` turns raw `.sql` files into C# classes with clean version tracking — making SQLite schema migrations easy to manage in Unity, or any project using embedded C# and SQLite.

> Inspired by [Goose](https://github.com/pressly/goose), built as an internal tool for Unity's ecosystem where no solid schema migration tooling exists but can be used for anything C#.

---

## ✨ Features

- CLI command to create new migrations with auto-incrementing numbers
- Converts `.up.sql`/`.down.sql` files into Unity-friendly `.cs` files
- Auto-generates a `MigrationRegistry.cs` file for easy runtime access
- Clean, no-runtime-dependency codegen
- Compatible with `sqlite-net` and other embedded SQLite libs

---

## 🧠 Why?

Using a sqlite database in a Unity project - although useful - seems niche
For that reason, there's virtually no (at least I couldn't find any) easy to use migration tools for a sqlite + Unity workflow
This tool fills that gap — letting you manage database evolution in a way that’s versioned, testable, and game-friendly.

---

## 📦 Install

### Option 1: Build from source

```bash
git clone https://github.com/mwac-dev/sqlinity
cd sqlinity
go build -o sqlinity
```

### Option 2: Make globally available (Windows/Linux/macOS)

Windows: Move sqlinity.exe somewhere like C:\Tools\sqlinity\ and add that folder to your System PATH

macOS/Linux: Move the binary to /usr/local/bin (or wherever you want and add to PATH)

Now you can run sqlinity from anywhere in your terminal.

---

## 🛠️ Usage

### 🗂 Project structure

```bash
my-unity-project/
├── config.json
├── migrations/
│   ├── 001_create_players.up.sql
│   ├── 001_create_players.down.sql
│   └── ...
├── Assets/
│   └── Scripts/
│       └── GeneratedMigrations/
```

### 📄 config.json

Create a config file in your project root:

```json 
{
  "sqlFolder": "./migrations",
  "outputFolder": "./Assets/Scripts/GeneratedMigrations",
  "namespace": "MyGame.Database.Migrations"
}
```

### 🚀 Commands

#### 🆕 Create a new migration

```bash
sqlinity create "create leaderboard table"
```

This will create:
```bash
migrations/003_create_leaderboard_table.up.sql
migrations/003_create_leaderboard_table.down.sql
```
Migration files will auto increment based on the existing files in the migrations directory.

#### 🛠 Generate C# migration classes

```bash
sqlinity generate
```

Outputs:
- Migration_003_create_leaderboard_table.cs
- MigrationRegistry.cs

Example MigrationRegistry.cs:
```csharp
public static class MigrationRegistry {
    public static readonly List<(string Name, string Sql)> Migrations = new() {
        (Migration_001_create_players.Name, Migration_001_create_players.Up),
        (Migration_002_add_score.Name, Migration_002_add_score.Up),
        (Migration_003_create_leaderboard_table.Name, Migration_003_create_leaderboard_table.Up)
    };
}
```

#### 🎮 How to Use in Unity
Use sqlite-net or similar with a startup script

```csharp
public static class SqliteMigrationRunner
{
    public static void RunMigrations(SQLiteConnection conn)
    {
        conn.Execute(@"
            CREATE TABLE IF NOT EXISTS __migrations (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
                name TEXT NOT NULL UNIQUE,
                applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
            );
        ");

        var applied = new HashSet<string>(
            conn.Query<string>("SELECT name FROM __migrations")
        );

        foreach (var (name, sql) in MigrationRegistry.Migrations)
        {
            if (applied.Contains(name)) continue;
            conn.Execute(sql);
            conn.Execute("INSERT INTO __migrations (name) VALUES (?)", name);
        }
    }
}
```
This ensures your user's local DB is always up to date with migrations when the app launches

---

# ✅ Naming Convention

File names must follow this structure:
```bash
NNN_descriptive_name.up.sql
NNN_descriptive_name.down.sql  (optional)
i.e.
001_create_players.up.sql
001_create_players.down.sql
```

---

# 💡 Good To Know

- .down.sql is optional — useful for dev resets but not so much at runtime (usually you want to ever migrate up/forward with any changes)
- Migrations are static .sql files — no parameters or dynamic values
- Generated C# files are safe to commit — they are your runtime format - just don't manually edit them, if you NEED to change them, edit the related .sql file and then re-generate using the tool

---

# 📌 Roadmap
- [x] Migration file parser

- [x] CLI create command

- [x] C# class generation for each migration

- [x] Auto-generated MigrationRegistry.cs

- [ ] Editor integration for Unity 

- [ ] --config flag support

---

# 🧑‍💻 License - MIT  
Use it in your games and tools if you want! Hopefully you find it useful for your use cases.
Tool comes with no warranties - as always use at your own risk and make sure to double check and test all sql queries locally before pushing a build to clients!
