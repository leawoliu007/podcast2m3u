package main

// content is stored in constant below
// var content embed.FS

const htmlContent = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Podcast2M3U Manager</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; max-width: 800px; margin: 0 auto; padding: 20px; background-color: #f4f4f9; }
        h1 { color: #333; }
        .card { background: white; padding: 20px; margin-bottom: 20px; border-radius: 8px; box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
        label { display: block; margin-bottom: 5px; font-weight: bold; }
        input[type="text"], input[type="url"] { width: 100%; padding: 8px; margin-bottom: 15px; border: 1px solid #ddd; border-radius: 4px; box-sizing: border-box; }
        button { background-color: #007bff; color: white; border: none; padding: 10px 15px; border-radius: 4px; cursor: pointer; }
        button:hover { background-color: #0056b3; }
        button.delete { background-color: #dc3545; }
        button.delete:hover { background-color: #c82333; }
        table { width: 100%; border-collapse: collapse; }
        th, td { text-align: left; padding: 12px; border-bottom: 1px solid #ddd; }
        tr:hover { background-color: #f1f1f1; }
        .status { margin-top: 10px; padding: 10px; border-radius: 4px; display: none; }
        .success { background-color: #d4edda; color: #155724; }
        .error { background-color: #f8d7da; color: #721c24; }
    </style>
</head>
<body>
    <h1>Podcast2M3U Manager</h1>

    <div class="card">
        <h2>Global Configuration</h2>
        <form id="globalConfigForm">
            <label>Update Interval (Cron)</label>
            <input type="text" id="updateInterval" name="update_interval" placeholder="0 * * * *">
            <label>Output Path</label>
            <input type="text" id="outputPath" name="output_path" placeholder="/path/to/playlists">
            <button type="submit">Save Global Config</button>
        </form>
    </div>

    <div class="card">
        <h2>Add Subscription</h2>
        <form id="addSubForm">
            <label>Name</label>
            <input type="text" id="subName" name="name" required placeholder="My Podcast">
            <label>RSS URL</label>
            <input type="url" id="subUrl" name="url" required placeholder="https://example.com/feed.xml">
            <label>Custom Schedule (Optional)</label>
            <input type="text" id="subCron" name="cron" placeholder="*/30 * * * *">
            <button type="submit">Add Subscription</button>
        </form>
    </div>

    <div class="card">
        <h2>Subscriptions</h2>
        <table id="subTable">
            <thead>
                <tr>
                    <th>Name</th>
                    <th>URL</th>
                    <th>Schedule</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody></tbody>
        </table>
    </div>

    <script>
        const API_BASE = '/api';

        async function loadConfig() {
            const res = await fetch(API_BASE + '/config');
            const data = await res.json();
            document.getElementById('updateInterval').value = data.global.update_interval || '';
            document.getElementById('outputPath').value = data.global.output_path || '';
            renderSubscriptions(data.subscriptions || []);
        }

        function renderSubscriptions(subs) {
            const tbody = document.querySelector('#subTable tbody');
            tbody.innerHTML = '';
            subs.forEach(sub => {
                const tr = document.createElement('tr');
                tr.innerHTML = '<td>' + sub.name + '</td>' +
                               '<td><a href="' + sub.url + '" target="_blank">Link</a></td>' +
                               '<td>' + (sub.cron || 'Default') + '</td>' +
                               '<td><button class="delete" onclick="deleteSub(\'' + sub.name + '\')">Delete</button></td>';
                tbody.appendChild(tr);
            });
        }

        document.getElementById('globalConfigForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const config = {
                global: {
                    update_interval: document.getElementById('updateInterval').value,
                    output_path: document.getElementById('outputPath').value
                }
            };
            await fetch(API_BASE + '/config', { method: 'POST', body: JSON.stringify(config) });
            loadConfig();
        });

        document.getElementById('addSubForm').addEventListener('submit', async (e) => {
            e.preventDefault();
            const sub = {
                name: document.getElementById('subName').value,
                url: document.getElementById('subUrl').value,
                cron: document.getElementById('subCron').value
            };
            await fetch(API_BASE + '/subscriptions', { method: 'POST', body: JSON.stringify(sub) });
            document.getElementById('addSubForm').reset();
            loadConfig();
        });

        async function deleteSub(name) {
            if (confirm('Are you sure you want to delete ' + name + '?')) {
                await fetch(API_BASE + '/subscriptions/' + encodeURIComponent(name), { method: 'DELETE' });
                loadConfig();
            }
        }

        loadConfig();
    </script>
</body>
</html>`
