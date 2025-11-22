import 'package:flutter/material.dart';
import 'dart:convert';
import '../env.dart';
import '../utils/api.dart' as api;

class ExercisesPage extends StatefulWidget {
  const ExercisesPage({super.key});

  @override
  State<ExercisesPage> createState() => _ExercisesPageState();
}

class _ExercisesPageState extends State<ExercisesPage> {
  final apiUrl = Env.instance.apiExerciseUrl;
  bool _isLoading = true;
  String? _error;
  List<dynamic> _items = [];

  @override
  void initState() {
    super.initState();
    _fetchItems();
  }

  Future<void> _fetchItems() async {
    try {
      final response = await api.sendProtectedGetRequest('$apiUrl/getForUser', {
        'limit': '50',
        'offset': '0',
      });
      setState(() {
        _items = jsonDecode(response['exercises']);
        _isLoading = false;
      });
    } catch (e) {
      setState(() {
        _error = 'Failed to load items';
        _isLoading = false;
      });
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return const Scaffold(body: Center(child: CircularProgressIndicator()));
    }
    if (_error != null) {
      return Scaffold(body: Center(child: Text(_error!)));
    }
    return Scaffold(
      appBar: AppBar(title: const Text('Exercises')),
      body: ListView.builder(
        itemCount: _items.length,
        itemBuilder: (context, index) {
          final item = _items[index];
          return ListTile(
            title: Text(item['title'] ?? ''),
            subtitle: Text(item['body'] ?? ''),
          );
        },
      ),
    );
  }
}
