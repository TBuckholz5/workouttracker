import 'dart:convert';
import 'dart:core';

import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:http/http.dart' as http;

final secureStorage = FlutterSecureStorage();

Future<Map<String, dynamic>> sendProtectedGetRequest(
  String url,
  Map<String, String> queryParams,
) async {
  String? token = await secureStorage.read(key: 'auth_token');
  if (token == null) {
    throw Exception('No auth token found in secure storage');
  }
  Map<String, String> headers = {
    "Authorization": "Bearer $token",
    "Content-Type": "application/json",
  };
  final uri = Uri.parse(url).replace(queryParameters: queryParams);
  final response = await http.get(uri, headers: headers);
  if (response.statusCode != 200) {
    throw Exception('Failed to get data: ${response.statusCode}');
  }
  return jsonDecode(response.body);
}

Future<Map<String, dynamic>> sendProtectedPostRequest(
  String url,
  Map<String, dynamic> payload,
) async {
  String? token = await secureStorage.read(key: 'auth_token');
  if (token == null) {
    throw Exception('No auth token found in secure storage');
  }
  Map<String, String> headers = {
    "Authorization": "Bearer $token",
    "Content-Type": "application/json",
  };
  final response = await http.post(
    Uri.parse(url),
    headers: headers,
    body: jsonEncode(payload),
  );
  if (response.statusCode != 200) {
    throw Exception('Failed to get data: ${response.statusCode}');
  }
  return jsonDecode(response.body);
}

Future<Map<String, dynamic>> sendPostRequest(
  String endpoint,
  Map<String, dynamic> payload,
) async {
  Map<String, String> headers = {'Content-Type': 'application/json'};
  final response = await http.post(
    Uri.parse(endpoint),
    headers: headers,
    body: jsonEncode(payload),
  );
  if (response.statusCode != 200) {
    throw Exception('Failed to post data: ${response.statusCode}');
  }
  return jsonDecode(response.body);
}
